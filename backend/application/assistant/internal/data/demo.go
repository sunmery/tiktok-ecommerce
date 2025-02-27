package data

import (
	"context"

	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"backend/api/order/v1"
	productv1 "backend/api/product/v1"
	"backend/application/assistant/internal/biz"


	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/go-kratos/kratos/v2/log"

	"google.golang.org/grpc"
)

type AssistantRepo struct {
	orderClient   orderv1.OrderServiceClient
	productClient productv1.ProductServiceClient
	log           *log.Helper
}

func NewAssistantRepo(
	orderConn *grpc.ClientConn,
	productConn *grpc.ClientConn,
	logger log.Logger,
) biz.OrderRepo {
	return &AssistantRepo{
		orderClient:   orderv1.NewOrderServiceClient(orderConn),
		productClient: productv1.NewProductServiceClient(productConn),
		log:           log.NewHelper(logger),
	}
}

func (r *AssistantRepo) CreateOrder(ctx context.Context, req *biz.OrderReq) (*orderv1.PlaceOrderResp, error) {
	// 1. 检查商品库存
	productResp, err := r.productClient.GetProduct(ctx, &productv1.GetProductRequest{
		Id:         req.Product,
		MerchantId: req.MerchantId.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("获取商品信息失败: %w", err)
	}

	if productResp.Stock <= 0 {
		return nil, biz.ErrProductOutOfStock
	}

	// 2. 构建订单请求
	orderReq := &orderv1.PlaceOrderReq{
		UserId:   req.UserID.String(),
		Currency: "CNY",
		Address: &orderv1.Address{
			StreetAddress: "AI Generated Address", // 需要根据业务补充真实地址
			City:          "Shanghai",
			State:         "Shanghai",
			Country:       "China",
			ZipCode:       "200000",
		},
		Email: "user@example.com",
		OrderItems: []*orderv1.OrderItem{
			{
				Item: &orderv1.CartItem{
					ProductId:  req.Product,
					Quantity:   1,
					MerchantId: productResp.MerchantId,
				},
				Cost: productResp.Price,
			},
		},
	}

	// 3. 调用订单服务
	resp, err := r.orderClient.PlaceOrder(ctx, orderReq)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	return resp, nil
}

func (r *AssistantRepo) QueryOrders(ctx context.Context, req *orderv1.ListOrderReq) (*orderv1.ListOrderResp, error) {
	resp, err := r.orderClient.ListOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	return resp, nil
}

// AI 生成响应实现
type AIGenerator struct {
	model *ark.ChatModel
}

func NewAIGenerator() *AIGenerator {
	return &AIGenerator{
		model: ark.NewChatModel(context.Background(), &ark.ChatModelConfig{
			Model:   "ep-20250205222056-l5n58",
			APIKey:  "9ca27fe9-6857-452b-a2ba-4b1713e46047",
			BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		}),
	}
}

func (g *AIGenerator) GenerateResponse(ctx context.Context, question string) (string, error) {
	template := prompt.FromMessages(
		schema.FString,
		&schema.Message{
			Role:    schema.System,
			Content: fmt.Sprintf(biz.SystemMessageTemplate, time.Now().Format(time.RFC3339)),
		},
		&schema.Message{
			Role:    schema.User,
			Content: fmt.Sprintf(biz.UserMessageTemplate, question),
		},
	)

	messages, err := template.Format(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("模板格式化失败: %w", err)
	}

	stream, err := g.model.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("模型调用失败: %w", err)
	}
	defer stream.Close()

	var sb strings.Builder
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("流式读取失败: %w", err)
		}
		sb.WriteString(resp.Content)
	}

	return sb.String(), nil
}

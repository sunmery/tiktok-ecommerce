CREATE TABLE "payment"."pay_record"
(
    "id"             int4                                        NOT NULL DEFAULT nextval('pay_record_id_seq'::regclass),
    "created_at"     timestamptz(6) NOT NULL DEFAULT now(),
    "deleted_at"     timestamptz(6) NOT NULL DEFAULT now(),
    "user_id"        varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "order_id"       varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "transcation_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "amount"         float8                                      NOT NULL,
    "pay_at"         timestamptz(6) NOT NULL,
    "status"         varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    CONSTRAINT "pay_record_pkey" PRIMARY KEY ("id")
)
;
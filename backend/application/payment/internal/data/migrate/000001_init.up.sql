CREATE SEQUENCE IF NOT EXISTS pay_record_id_seq;

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

ALTER TABLE "public"."pay_record"
    OWNER TO "postgres";

CREATE INDEX "idx_order_id" ON "public"."pay_record" USING btree (
    "order_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );

CREATE INDEX "idx_transcation_id" ON "public"."pay_record" USING btree (
    "transcation_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );

CREATE INDEX "idx_user_id" ON "public"."pay_record" USING btree (
    "user_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );
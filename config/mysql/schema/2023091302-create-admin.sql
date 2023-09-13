CREATE SCHEMA IF NOT EXISTS `furumane` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `furumane`.`admins` (
  `id`            VARCHAR(22)  NOT NULL,          -- ユーザーID
  `cognito_id`    VARCHAR(36)  NOT NULL,          -- ユーザーID（Cognito用）
  `provider_type` INT          NOT NULL,          -- 認証種別
  `email`         VARCHAR(256) NULL DEFAULT NULL, -- メールアドレス
  `exists`        TINYINT      NULL DEFAULT 1,    -- 有効化フラグ
  `created_at`    DATETIME     NOT NULL,          -- 登録日時
  `updated_at`    DATETIME     NOT NULL,          -- 更新日時
  `verified_at`   DATETIME     NULL DEFAULT NULL, -- 確認日時
  `deleted_at`    DATETIME     NULL DEFAULT NULL, -- 退会日時
  PRIMARY KEY(`id`)
);

CREATE UNIQUE INDEX `ui_admin_cognito_id` ON `furumane`.`admins` (`cognito_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_admin_email` ON `furumane`.`admins` (`exists` DESC, `email` ASC) VISIBLE;

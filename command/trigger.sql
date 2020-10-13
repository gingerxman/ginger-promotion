USE ginger_mall;

DROP TRIGGER IF EXISTS order_has_product_insert_trigger;
DELIMITER $$
CREATE TRIGGER order_has_product_insert_trigger AFTER INSERT ON ginger_mall.order_has_product FOR EACH ROW
BEGIN
    UPDATE product_pool_product SET sold_count=sold_count+new.count
    WHERE id=new.pool_product_id;
END
$$
DELIMITER ;

DROP TRIGGER IF EXISTS order_insert_trigger;
DELIMITER $$
CREATE TRIGGER order_insert_trigger AFTER INSERT ON ginger_mall.order_order FOR EACH ROW
BEGIN
    DECLARE record_count TINYINT;

    SELECT count(*) into record_count FROM order_user_consumption_record WHERE user_id = new.user_id AND corp_id = new.corp_id;

    IF record_count = 0 THEN
        INSERT INTO order_user_consumption_record (user_id, corp_id, money, consume_count, created_at, updated_at)
        VALUES (new.user_id, new.corp_id, new.final_money, 1, now(), now());
    ELSE
        UPDATE order_user_consumption_record SET money=money+new.final_money, consume_count=consume_count+1, updated_at=now();
    END IF;
END
$$
DELIMITER ;

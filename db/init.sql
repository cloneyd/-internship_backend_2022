DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS operations;

CREATE TABLE balances
(
    user_id         INT PRIMARY KEY,
    balance_amount  NUMERIC,

    CONSTRAINT positive_balance
        CHECK (balance_amount > 0.0)
);

CREATE TABLE reserve_balances
(
    user_id         INT PRIMARY KEY,
    reserve_amount  NUMERIC,

    CONSTRAINT positive_balance
        CHECK (reserve_amount > 0.0)
);

CREATE TABLE services
(
    service_id          SERIAL PRIMARY KEY,
    service_title       VARCHAR(128),
    service_description TEXT
);

CREATE TABLE operations
(
    operation_id    SERIAL PRIMARY KEY,
    order_id        INT NOT NULL,
    user_id         INT NOT NULL,
    service_id      INT NOT NULL,
    price           NUMERIC NOT NULL,
    approved        BOOLEAN,
    done            BOOLEAN NOT NULL,
    done_timestamp  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES balances (user_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE,

    CONSTRAINT fk_service_id
        FOREIGN KEY (service_id)
            REFERENCES services (service_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
);


CREATE
OR REPLACE PROCEDURE deposit_on_balance(user_id_ INT, amount_ NUMERIC)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO balances (user_id, balance_amount)
    VALUES (user_id_, amount_)
    ON CONFLICT (user_id)
    DO UPDATE
        SET balance_amount = balances.balance_amount + EXCLUDED.balance_amount;
    COMMIT;
END;
$$;

CREATE
OR REPLACE PROCEDURE reserve_amount(user_id_ INT, service_id_ INT, order_id_ INT, price_ NUMERIC)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO operations (order_id, user_id, service_id, price, done)
    VALUES (order_id_, user_id_, service_id_, price_, false);

    IF NOT EXISTS (SELECT FROM balances AS b
                   WHERE b.user_id = user_id AND
                   b.balance_amount > price_) THEN
        ROLLBACK;
    ELSE
        INSERT INTO reserve_balances (user_id, reserve_amount)
        VALUES (user_id_, price_)
        ON CONFLICT (user_id) DO NOTHING;

        UPDATE balances AS b
        SET balance_amount = balance_amount - price_
        WHERE b.user_id = user_id_;

        UPDATE reserve_balances AS rb
        SET reserve_amount = reserve_amount + price_
        WHERE rb.user_id = user_id_;

        COMMIT;
    END IF;
END;
$$;

CREATE
OR REPLACE PROCEDURE approve_order(user_id_ INT, service_id_ INT, order_id_ INT, price_ NUMERIC)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE operations AS o
    SET approved = true, done = true, done_timestamp = CURRENT_TIMESTAMP
    WHERE o.order_id   = order_id_ AND
          o.service_id = service_id_ AND
          o.user_id    = user_id_ AND
          o.price      = price_ AND
          !done;

    IF NOT EXISTS (SELECT FROM reserve_balances AS rb
                   WHERE rb.user_id = user_id_ AND
                   rb.reserve_amount > price_) THEN
        ROLLBACK;
    ELSE
        UPDATE reserve_balances AS rb
        SET reserve_amount = reserve_amount - price_
        WHERE rb.user_id = user_id_;

        COMMIT;
    END IF;
END;
$$;

CREATE
OR REPLACE PROCEDURE disapprove_order(user_id_ INT, service_id_ INT, order_id_ INT, price_ NUMERIC)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE operations AS o
    SET approved = false, done = true, done_timestamp = CURRENT_TIMESTAMP
    WHERE o.order_id   = order_id_ AND
          o.service_id = service_id_ AND
          o.user_id    = user_id_ AND
          o.price      = price_;

    UPDATE balances as b
    SET balance_amount = balance_amount + price_
    WHERE b.user_id = user_id_;

    UPDATE reserve_balances AS rb
    SET reserve_amount = reserve_amount - price_
    WHERE rb.user_id = user_id_;

    COMMIT;
END;
$$;

CREATE
OR REPLACE FUNCTION get_service_month_revenue_report(report_period DATE)
    RETURNS TABLE (
        service_name VARCHAR(128),
        revenue      NUMERIC
    )
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY SELECT
        s.service_title as "service_title",
        SUM(o.price) as "revenue"
    FROM operations AS o
    INNER JOIN
        services AS s ON s.service_id = o.service_id
    WHERE EXTRACT(MONTH FROM o.done_timestamp) = EXTRACT(MONTH FROM report_period) AND
          EXTRACT(YEAR FROM o.done_timestamp)  = EXTRACT(YEAR FROM report_period) AND
          o.approved
    GROUP BY s.service_title;
END;
$$;

-- SELECT get_service_month_revenue_report('11-01-2022');

INSERT INTO services (service_title, service_description)
VALUES ('car buy', 'Mclaren F1'),
       ('home rent', 'home in LA');

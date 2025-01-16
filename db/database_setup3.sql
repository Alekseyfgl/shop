-- =========================================
-- 1. Создать схему shop
-- =========================================
CREATE SCHEMA shop;

-- =========================================
-- 2. Таблица characteristics
--    (характеристики товаров/объектов)
-- =========================================
CREATE TABLE shop.characteristics
(
    id             SERIAL PRIMARY KEY,
    title          TEXT NOT NULL,
    description    TEXT,
    is_visible BOOL DEFAULT TRUE
);

-- =========================================
-- 3. Таблица node_types
--    (типы "узлов": например, товар, категория, услуга и т.п.)
-- =========================================
CREATE TABLE shop.node_types
(
    id          SERIAL PRIMARY KEY,
    type        TEXT NOT NULL,
    description TEXT
);

-- =========================================
-- 4. Таблица nodes
--    (сами "узлы": товары, категории и т.п.)
-- =========================================
CREATE TABLE shop.nodes
(
    id           SERIAL PRIMARY KEY,
    title        TEXT      NOT NULL,
    node_type_id INT REFERENCES shop.node_types (id)
        ON DELETE CASCADE,
    description  TEXT,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    removed_at   TIMESTAMP          DEFAULT NULL
);

-- Индекс на created_at (если нужно часто искать по дате создания)
CREATE INDEX idx_nodes_created_at
    ON shop.nodes (created_at);

-- =========================================
-- 5. Таблица units_measurement
--    (справочник единиц измерения)
-- =========================================
CREATE TABLE shop.units_measurement
(
    id          SERIAL PRIMARY KEY,
    unit        VARCHAR(10) NOT NULL, -- напр. "шт", "кг", "литр"
    description TEXT
);

-- =========================================
-- 6. Таблица quantity
--    (содержит связь единицы измерения + количество)
-- =========================================
CREATE TABLE shop.quantity
(
    id                  SERIAL PRIMARY KEY,
    unit_measurement_id INT NOT NULL REFERENCES shop.units_measurement (id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    quantity            INT NOT NULL
);

CREATE TABLE shop.char_default_value
(
    id                SERIAL PRIMARY KEY,
    characteristic_id INT NOT NULL REFERENCES shop.characteristics (id) ON DELETE CASCADE ON UPDATE CASCADE,
    value             TEXT,
    UNIQUE (characteristic_id, value)
);

-- =========================================
-- 7. Таблица characteristic_values
--    (значения характеристик для каждого "узла" / товара)
-- =========================================
CREATE TABLE shop.characteristic_values
(
    node_id           INT  NOT NULL
        REFERENCES shop.nodes (id)
            ON DELETE CASCADE,
    characteristic_id INT  NOT NULL
        REFERENCES shop.characteristics (id)
            ON DELETE CASCADE,
    add_params        jsonb,
    value             TEXT NOT NULL,

    -- Делаем уникальную комбинацию для бизнес-логики,
    -- если нужно гарантировать, что не будет повторов:
    UNIQUE (node_id, characteristic_id, value)
);

-- Индексы для ускорения поиска
CREATE INDEX idx_characteristic_values_node_id
    ON shop.characteristic_values (node_id);

CREATE INDEX idx_characteristic_values_characteristic_id
    ON shop.characteristic_values (characteristic_id);

CREATE INDEX idx_characteristic_values_value
    ON shop.characteristic_values (value);


CREATE OR REPLACE FUNCTION set_created_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.created_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_created_at
    BEFORE INSERT
    ON shop.nodes
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();


CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Check if the update affects any column other than `removed_at`
    IF NOT TG_OP = 'UPDATE' OR ROW (NEW.*) IS DISTINCT FROM ROW (OLD.*) THEN
        NEW.updated_at := NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_updated_at
    BEFORE UPDATE
    ON shop.nodes
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- ======================================

CREATE TABLE size
(
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL UNIQUE,
    description TEXT
);
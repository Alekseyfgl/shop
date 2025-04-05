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
    removed_at   TIMESTAMP          DEFAULT NULL,
    images TEXT DEFAULT NULL,
    price_byn INT,
    price_run INT
);

-- Индекс на created_at (если нужно часто искать по дате создания)
CREATE INDEX idx_nodes_created_at
    ON shop.nodes (created_at);



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



-- ======================================
-- ======================================
-- SEARCH
-- ======================================
-- ======================================

ALTER TABLE shop.nodes
    ADD COLUMN search_vector tsvector;

CREATE INDEX idx_nodes_search_vector
    ON shop.nodes USING GIN (search_vector);


CREATE OR REPLACE FUNCTION fn_update_nodes_search_vector()
    RETURNS TRIGGER AS
$$
DECLARE
    v_type        text;
    v_char_values text;
BEGIN
    -- 1) Получаем тип из node_types
    SELECT nt.type
    INTO v_type
    FROM shop.node_types nt
    WHERE nt.id = NEW.node_type_id;

    -- 2) Агрегируем только "видимые" характеристики:
    --    - для каждой строки получаем c.name и cv.value
    --    - соединяем их через пробел или через двоеточие "name: value"
    --    - итоговые пары объединяем в одну строку через пробел
    SELECT string_agg(
                   c.title || ' ' || cv.value -- или c.name || ': ' || cv.value
               , ' '
           )
    INTO v_char_values
    FROM shop.characteristic_values cv
             JOIN shop.characteristics c ON c.id = cv.characteristic_id
    WHERE cv.node_id = NEW.id
      AND c.is_visible = true;

    -- 3) Собираем финальную строку для to_tsvector:
    --    title + description + type + (название_хар + значение)
    NEW.search_vector := to_tsvector(
            'russian', -- можно поменять на 'simple', 'english' и т. п.
            COALESCE(NEW.title, '') || ' '
                || COALESCE(NEW.description, '') || ' '
                || COALESCE(v_type, '') || ' '
                || COALESCE(v_char_values, '')
                         );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;



CREATE TRIGGER tr_update_nodes_search_vector
    BEFORE INSERT OR UPDATE
    ON shop.nodes
    FOR EACH ROW
EXECUTE PROCEDURE fn_update_nodes_search_vector();


CREATE OR REPLACE FUNCTION fn_characteristics_update_nodes()
    RETURNS TRIGGER AS
$$
BEGIN
    -- Триггер вызывает обновление nodes,
    -- чтобы заставить сработать триггер на nodes (tr_update_nodes_search_vector).
    UPDATE shop.nodes
    SET title = title -- любое "псевдо-обновление", лишь бы триггер на nodes сработал
    WHERE id = NEW.node_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER tr_char_values_insert
    AFTER INSERT
    ON shop.characteristic_values
    FOR EACH ROW
EXECUTE PROCEDURE fn_characteristics_update_nodes();

CREATE TRIGGER tr_char_values_update
    AFTER UPDATE
    ON shop.characteristic_values
    FOR EACH ROW
EXECUTE PROCEDURE fn_characteristics_update_nodes();

CREATE TRIGGER tr_char_values_delete
    AFTER DELETE
    ON shop.characteristic_values
    FOR EACH ROW
EXECUTE PROCEDURE fn_characteristics_update_nodes();




SELECT *
FROM shop.nodes
WHERE search_vector @@ plainto_tsquery('russian', 'майка летняя XXL');


-- обновить search_vector в nodes где есть null
UPDATE shop.nodes n
SET search_vector = to_tsvector(
        'russian', -- Или 'simple', 'english' и т.п., в зависимости от языка
        COALESCE(n.title, '') || ' '
            || COALESCE(n.description, '') || ' '
            || COALESCE(nt.type, '') || ' '
            || COALESCE(
                (SELECT string_agg(
                                c.title || ' ' || cv.value,
                                ' '
                        )
                 FROM shop.characteristic_values cv
                          JOIN shop.characteristics c ON c.id = cv.characteristic_id
                 WHERE cv.node_id = n.id
                   AND c.is_visible = true),
                ''
               )
                    )
FROM shop.node_types nt
WHERE n.search_vector IS NULL
  AND nt.id = n.node_type_id;




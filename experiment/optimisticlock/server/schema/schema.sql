CREATE TABLE coupon (
    "id" SERIAL,
    "code" varchar(50),
    "status" int,
    PRIMARY KEY("id")
);
--end

CREATE TABLE coupon_quantity (
    "id" SERIAL,
    "coupon_id" int,
    "quantity" int,
    PRIMARY KEY("id")
);
--end

CREATE TABLE coupon_history (
    "id" SERIAL,
    "coupon_id" int,
    "modified_quantity" int,
    PRIMARY KEY("id")
);   
--end
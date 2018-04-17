CREATE TABLE IF  NOT EXISTS gd(
  `code`  VARCHAR (30) NOT NULL,
	`s_name`   VARCHAR (30) NOT NULL,
	`risefall_rate`   VARCHAR (30) NOT NULL,
	`change_rate`   VARCHAR (30) NOT NULL,
	`lastest_price`  VARCHAR (30) NOT NULL,
	`previous_points`   VARCHAR (30) NOT NULL,
	`create_time` DATE DEFAULT NULL COMMENT '创建时间',
	`previous_points_date` DATE ,
	`data_types`      INT(4) NOT NULL COMMENT '1.创新高 2.创新低'
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
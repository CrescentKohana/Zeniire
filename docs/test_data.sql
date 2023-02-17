-- INSERT INTO records (uuid, amount, datetime) VALUES (
--                                                      gen_random_uuid(),
--                                                      trunc(random()*100000000),
--                                                      random() * (timestamp '2020-01-01 12:00:00' - timestamp '2023-01-01 12:00:00')
--                                                      );

INSERT INTO records (uuid, amount, datetime)
SELECT gen_random_uuid(),
       (random()*100000000)::integer,
       DATE '2020-01-01' + (random() * 700)::integer
FROM generate_series(1, 1000000);

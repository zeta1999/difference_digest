CREATE OR REPLACE FUNCTION f_hash(idx int, element bigint) RETURNS numeric
    RETURNS NULL ON NULL INPUT
    IMMUTABLE
    LANGUAGE plpythonu
AS $$
  hashSeeds = [18269989962351869307, 9143901319630896501, 2072764263930962169, 417226483919003799, 16485935163296413021]
  mask = int("0x5555555555555555", 16)

  inner = (mask * (element ^ (element >> 32))) % 18446744073709551616
  return (hashSeeds[idx] * inner) % 18446744073709551616
$$; 

CREATE OR REPLACE FUNCTION f_trailing_zeros(element numeric) RETURNS int
    RETURNS NULL ON NULL INPUT
    IMMUTABLE
    LANGUAGE plpythonu
AS $$
  x = int(element)
  count = 0
  if x == 0:
    return count

  while ((x & 1) == 0): 
      x = x >> 1
      count += 1
      
  return count 
$$; 

CREATE OR REPLACE FUNCTION bit_xor_sfunc(agg bigint, value bigint) RETURNS bigint
  RETURNS NULL ON NULL INPUT
  IMMUTABLE
  LANGUAGE SQL
AS $$
SELECT agg # value;
$$;

CREATE OR REPLACE AGGREGATE bit_xor (bigint) (
    SFUNC = bit_xor_sfunc,
    STYPE = bigint,
    INITCOND = 0
);

CREATE OR REPLACE FUNCTION bit_xor_numeric_sfunc(agg numeric, value numeric) RETURNS numeric
  RETURNS NULL ON NULL INPUT
  IMMUTABLE
  LANGUAGE plpythonu
AS $$
  return int(agg) ^ int(value)
$$;

CREATE OR REPLACE AGGREGATE bit_xor_numeric (numeric) (
    SFUNC = bit_xor_numeric_sfunc,
    STYPE = numeric,
    INITCOND = 0
);



/*
SELECT 
  f_trailing_zeros(f_hash(3 + 1, id)) AS estimator, 
  f_hash(idx, id) % 80 AS cell, 
  bit_xor(id::bigint) AS id_sum, 
  bit_xor_numeric(f_hash(3 + 0, id)) AS hash_sum,  
  COUNT(id) AS count
FROM (
    SELECT 0 AS idx, * FROM mythings UNION SELECT 1, * FROM mythings UNION SELECT 2, * FROM mythings
  ) things 
GROUP BY 1, 2 
HAVING  f_trailing_zeros(f_hash(3 + 1, id)) = 0
ORDER BY 1, 2;

SELECT 
  id, idx,
  f_trailing_zeros(f_hash(3 + 1, id)) AS estimator, 
  f_hash(idx, id) % 80 AS cell
FROM (
    SELECT 0 AS idx, * FROM mythings UNION SELECT 1, * FROM mythings UNION SELECT 2, * FROM mythings
  ) things 
WHERE id < 5
ORDER BY 1, 2;

SELECT i, f_hash(0, i) FROM generate_series(1,20) AS s(i);
SELECT  f_hash(3 + 1, id), f_trailing_zeros(f_hash(3 + 1, id)) AS estimator FROM generate_series(0,12) AS s(id);
*/

-- Insert sale offers with user_id from 2 to 10
INSERT INTO sale_offers (user_id, description, price, date_of_issue, margin, status, is_auction) VALUES
-- Regular sales (non-auction)
(2, 'Used compact sedan 2018, low mileage, single owner, excellent condition', 65000, '2025-01-10 08:30:00', 5, 'pending', FALSE),
(2, 'Compact sedan 2019, dark exterior, 45K km, leather interior, rearview camera', 72000, '2025-01-15 10:45:00', 3, 'pending', FALSE),
(3, 'Mid-size sedan 2020, premium trim, complete maintenance history', 120000, '2025-02-01 14:20:00', 10, 'pending', FALSE),
(3, 'Compact executive sedan 2017, sport trim, well maintained', 95000, '2025-01-20 09:15:00', 5, 'pending', FALSE),
(4, 'Sport compact hatchback 2019, manual transmission, tuned performance', 85000, '2025-01-05 16:30:00', 5, 'pending', FALSE),
(4, 'Luxury sedan 2016, performance package, panoramic sunroof', 88000, '2025-02-10 11:00:00', 10, 'pending', FALSE),
(5, 'Convertible roadster 2021, only 15K km', 110000, '2025-01-25 13:45:00', 3, 'pending', FALSE),
(5, 'Muscle coupe 2018, V8 engine, premium audio system', 130000, '2025-02-05 15:20:00', 10, 'pending', FALSE),
(6, 'Compact SUV 2020, full options, like new', 92000, '2025-01-18 09:50:00', 5, 'pending', FALSE),
(6, 'Compact SUV 2019, diesel engine, economical, family-oriented', 78000, '2025-02-12 14:15:00', 3, 'pending', FALSE),
(7, 'Station wagon 2017, spacious interior', 52000, '2025-01-22 10:30:00', 5, 'pending', FALSE),
(7, 'Compact hatchback 2020, fuel efficient', 65000, '2025-02-08 12:45:00', 3, 'pending', FALSE),
(8, 'Electric sedan 2021, advanced driver-assist, premium connectivity', 195000, '2025-01-30 16:20:00', 10, 'pending', FALSE),
(8, 'Compact electric hatchback 2019, extended range', 89000, '2025-02-15 09:30:00', 5, 'pending', FALSE),
(9, 'Compact sedan 2018, diesel engine, large cargo capacity', 68000, '2025-01-28 14:00:00', 3, 'pending', FALSE),
(9, 'Compact SUV 2020, panoramic sunroof', 97000, '2025-02-03 11:50:00', 5, 'pending', FALSE),
(10, 'City car 2019, low fuel consumption', 45000, '2025-01-12 15:15:00', 3, 'pending', FALSE),
(10, 'Mid-size SUV 2017, advanced safety features', 115000, '2025-02-07 13:40:00', 10, 'pending', FALSE),
(2, 'Compact hatchback 2020, sporty design, well equipped', 75000, '2025-01-17 08:55:00', 5, 'pending', FALSE),
(3, 'Compact sedan 2018, comfortable ride, good condition', 58000, '2025-02-09 10:10:00', 3, 'pending', FALSE),
(4, 'Compact hybrid SUV 2019, spacious interior', 110000, '2025-01-23 14:35:00', 10, 'pending', FALSE),
(5, 'Compact hatchback 2020, sporty styling', 72000, '2025-02-11 09:25:00', 5, 'pending', FALSE),
(6, 'Compact crossover 2018, all-wheel drive, off-road capable', 89000, '2025-01-19 11:15:00', 5, 'pending', FALSE),
(7, 'Subcompact hatchback 2021, stylish, ideal for city driving', 95000, '2025-02-02 16:00:00', 3, 'pending', FALSE),
(8, 'Budget-friendly SUV 2019, economical operation', 62000, '2025-01-26 12:30:00', 5, 'pending', FALSE),

-- Auctions
(2, 'Classic sports coupe 1985, collector item, vintage condition', 250000, '2025-01-05 09:00:00', 10, 'pending', FALSE),
(3, 'Rare classic roadster 1990, exceptional condition', 180000, '2025-01-10 14:30:00', 10, 'pending', FALSE),
(4, 'Exotic sports coupe 1997, perfect mechanical condition', 320000, '2025-01-15 16:45:00', 10, 'pending', FALSE),
(5, 'Vintage coupe 1968, restored, original components', 230000, '2025-01-20 11:20:00', 10, 'pending', FALSE),
(6, 'Exotic two-door 2008, low mileage, luxury finish', 450000, '2025-01-25 15:10:00', 10, 'pending', FALSE),
(7, 'Classic British sports car 1964, collector’s edition', 550000, '2025-01-30 13:40:00', 10, 'pending', FALSE),
(8, 'Luxury grand tourer 2010, exceptional condition', 390000, '2025-02-04 10:15:00', 10, 'pending', FALSE),
(9, 'Classic American sports car 1975, iconic design', 180000, '2025-02-09 14:00:00', 10, 'pending', FALSE),
(10, 'Luxury grand tourer 2012, high-performance', 280000, '2025-02-14 12:30:00', 10, 'pending', FALSE),
(2, 'Classic off-road 2000, rugged design', 150000, '2025-01-08 09:45:00', 5, 'pending', FALSE),
(3, 'Performance coupe 2012, track-ready, modified', 190000, '2025-01-13 16:20:00', 5, 'pending', FALSE),
(4, 'Muscle coupe 2015, powerful engine', 170000, '2025-01-18 11:50:00', 5, 'pending', FALSE),
(5, 'Italian sports coupe 2016, lightweight design', 210000, '2025-01-23 14:15:00', 5, 'pending', FALSE),
(6, 'High-performance rally-inspired coupe 2014', 130000, '2025-01-28 10:30:00', 5, 'pending', FALSE),
(7, 'Lightweight sports roadster 2010, minimalistic design', 160000, '2025-02-02 15:45:00', 5, 'pending', FALSE),
(8, 'Elegant grand tourer 2013, luxury amenities', 240000, '2025-02-07 12:00:00', 5, 'pending', FALSE),
(9, 'Luxury performance sedan 2016, high output', 185000, '2025-02-12 09:20:00', 5, 'pending', FALSE),
(10, 'Rare supercar 2011, limited production', 950000, '2025-02-17 16:40:00', 10, 'pending', FALSE),
(2, 'Mid-engine sports coupe 2017, dynamic handling', 220000, '2025-01-11 13:10:00', 5, 'pending', FALSE),
(3, 'Everyday supercar 2014, road-legal performance', 380000, '2025-01-16 10:50:00', 10, 'pending', FALSE),
(4, 'High-performance sports coupe 2015, advanced engineering', 290000, '2025-01-21 15:30:00', 5, 'pending', FALSE),
(5, 'Exotic supercar 2018, cutting-edge performance', 650000, '2025-01-26 12:45:00', 10, 'pending', FALSE),
(6, 'Legendary hypercar 2008, top-tier performance', 2500000, '2025-01-31 14:55:00', 10, 'pending', FALSE);


INSERT INTO auctions (offer_id, date_end, buy_now_price, initial_price) VALUES
  (26, '2025-06-11 12:00:00+00', 275000, 250000),
  (27, '2025-06-11 12:00:00+00', 198000, 180000),
  (28, '2025-06-11 12:00:00+00', 352000, 320000),
  (29, '2025-06-11 12:00:00+00', 253000, 230000),
  (30, '2025-06-11 12:00:00+00', 495000, 450000),
  (31, '2025-06-11 12:00:00+00', 605000, 550000),
  (32, '2025-06-11 12:00:00+00', 429000, 390000),
  (33, '2025-06-11 12:00:00+00', 198000, 180000),
  (34, '2025-06-11 12:00:00+00', 308000, 280000),
  (35, '2025-06-11 12:00:00+00', 165000, 150000),
  (36, '2025-06-11 12:00:00+00', 209000, 190000),
  (37, '2025-06-11 12:00:00+00', 187000, 170000),
  (38, '2025-06-12 12:00:00+00', 231000, 210000),
  (39, '2025-06-12 12:00:00+00', 143000, 130000),
  (40, '2025-06-12 12:00:00+00', 176000, 160000),
  (41, '2025-06-12 12:00:00+00', 264000, 240000),
  (42, '2025-06-12 12:00:00+00', 203500, 185000),
  (43, '2025-06-12 12:00:00+00', 1045000, 950000),
  (44, '2025-06-12 12:00:00+00',    NULL, 220000),
  (45, '2025-06-12 12:00:00+00', 418000, 380000),
  (46, '2025-06-12 12:00:00+00', 319000, 290000),
  (47, '2025-06-12 12:00:00+00',    NULL, 650000),  
  (48, '2025-06-12 12:00:00+00', 2750000, 2500000);



INSERT INTO cars (
  offer_id,
  vin,
  production_year,
  mileage,
  number_of_doors,
  number_of_seats,
  engine_power,
  engine_capacity,
  registration_number,
  registration_date,
  color,
  fuel_type,
  drive,
  transmission,
  number_of_gears,
  model_id
) VALUES
  -- Porsche 911 (1985) auction
  (26,
   'WP0AB2A90FS123456',
   1985,
   125000,
   2,
   2,
   285,
   2994,
   'P911CL85',
   '1985-06-15',
   'red',
   'petrol',
   'rwd',
   'manual',
   5,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Porsche'
       AND m.name = '911'        -- model '911' :contentReference[oaicite:0]{index=0}
   )
  ),
  -- Mercedes-Benz 300SL (1990) auction
  (27,
   'WDBBA48D1LA123456',
   1990,
   98000,
   2,
   2,
   228,
   2996,
   'MB300SL90',
   '1990-04-20',
   'black',
   'petrol',
   'rwd',
   'manual',
   4,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Mercedes-Benz'
       AND m.name = '300SL'     -- model '300SL' :contentReference[oaicite:1]{index=1}
   )
  ),
  -- Ferrari F355 (1997) auction
  (28,
   'ZFFZS17A7V0101234',
   1997,
   43000,
   2,
   2,
   380,
   3496,
   'F355SP97',
   '1997-08-05',
   'yellow',
   'petrol',
   'rwd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Ferrari'
       AND m.name = 'F355/355 F1'  -- model 'F355/355 F1' :contentReference[oaicite:2]{index=2}
   )
  ),
  -- Ford Mustang (1968) auction
  (29,
   '7R04S123456',
   1968,
   75000,
   2,
   4,
   290,
   4730,
   'MSTNG68',
   '1968-05-22',
   'blue',
   'petrol',
   'rwd',
   'manual',
   4,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Ford'
       AND m.name = 'Mustang'    -- model 'Mustang' :contentReference[oaicite:3]{index=3}
   )
  ),
  -- Lamborghini Gallardo (2008) auction
  (30,
   'ZHWGU11P08LA12345',
   2008,
   15000,
   2,
   2,
   500,
   5204,
   'GALL2008',
   '2008-03-10',
   'orange',
   'petrol',
   'rwd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Lamborghini'
       AND m.name = 'Gallardo Coupe'  -- model 'Gallardo Coupe' :contentReference[oaicite:4]{index=4}
   )
  ),
  -- Jaguar E-Type (1964) auction
  (31,
   'SWRJB20RXWOC12345',
   1964,
   89000,
   2,
   2,
   265,
   3781,
   'JETYPE64',
   '1964-07-14',
   'green',
   'petrol',
   'rwd',
   'manual',
   4,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Jaguar'
       AND m.name = 'E-Type'    -- assuming E-Type is in your Jaguar models
   )
  ),
  -- Aston Martin DB9 (2010) auction
  (32,
   'SCFAB02D28GK12345',
   2010,
   30000,
   2,
   4,
   470,
   5935,
   'ASTDB910',
   '2010-05-30',
   'black',
   'petrol',
   'rwd',
   'automatic',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Aston Martin'
       AND m.name = 'DB9'        -- model 'DB9'
   )
  ),
  -- Chevrolet Corvette (1975) auction
  (33,
   '1Z37L5S5135801234',
   1975,
   102000,
   2,
   2,
   350,
   5700,
   'CORV1975',
   '1975-09-12',
   'red',
   'petrol',
   'rwd',
   'manual',
   4,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Chevrolet'
       AND m.name = 'Corvette'  -- model 'Corvette'
   )
  ),
  -- Bentley Continental GT (2012) auction
  (34,
   'SCBBB7ZA2CC123456',
   2012,
   25000,
   2,
   4,
   552,
   5998,
   'BENTGT12',
   '2012-03-21',
   'brown',
   'petrol',
   'awd',
   'automatic',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Bentley'
       AND m.name = 'Continental GT'  -- model 'Continental GT'
   )
  ),
  -- Land Rover Defender (2000) auction
  (35,
   'SALDV1245YA123456',
   2000,
   140000,
   4,
   5,
   174,
   2287,
   'DEF2000',
   '2000-08-01',
   'green',
   'diesel',
   'awd',
   'manual',
   5,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Land Rover'
       AND m.name = 'Defender 110'   -- model 'Defender 110'
   )
  ),
  -- BMW M3 (2012) auction
  (36,
   'WBSBL9C54CE112345',
   2012,
   40000,
   2,
   4,
   414,
   4395,
   'BMWM312',
   '2012-07-10',
   'gray',
   'petrol',
   'rwd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'BMW'
       AND m.name = 'M3'         -- model 'M3'
   )
  ),
  -- Dodge Challenger (2015) auction
  (37,
   '2C3CDZBT5FH123456',
   2015,
   22000,
   2,
   5,
   375,
   3600,
   'DODGCH15',
   '2015-06-01',
   'black',
   'petrol',
   'rwd',
   'automatic',
   8,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Dodge'
       AND m.name = 'Challenger' -- model 'Challenger'
   )
  ),
  -- Alfa Romeo 4C (2016) auction
  (38,
   'ZARDE52V3G1234567',
   2016,
   15000,
   2,
   2,
   237,
   1742,
   'ALF4C16',
   '2016-05-20',
   'red',
   'petrol',
   'rwd',
   'dual_clutch',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Alfa Romeo'
       AND m.name = '4C'         -- model '4C'
   )
  ),
  -- Subaru Impreza WRX STI (2014) auction
  (39,
   'JF1GV8JS5EL123456',
   2014,
   30000,
   4,
   5,
   305,
   2457,
   'SUBUWRX14',
   '2014-06-12',
   'blue',
   'petrol',
   'awd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Subaru'
       AND m.name = 'Impreza'    -- model 'Impreza'
   )
  ),
  -- Lotus Elise (2010) auction
  (40,
   'SCCPC111AC1234567',
   2010,
   18000,
   2,
   2,
   190,
   1796,
   'LOTUEL10',
   '2010-07-05',
   'yellow',
   'petrol',
   'rwd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Lotus'
       AND m.name = 'Elise'      -- model 'Elise'
   )
  ),
  -- Maserati GranTurismo (2013) auction
  (41,
   'ZAM45VMA6D1234567',
   2013,
   25000,
   2,
   4,
   454,
   4699,
   'MASGTC13',
   '2013-08-18',
   'white',
   'petrol',
   'rwd',
   'automatic',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Maserati'
       AND m.name = 'GranTurismo' -- model 'GranTurismo'
   )
  ),
  -- Cadillac CTS-V (2016) auction
  (42,
   '1GYS4HKJ2GR123456',
   2016,
   28000,
   4,
   5,
   640,
   6162,
   'CADCTS16',
   '2016-06-25',
   'black',
   'petrol',
   'rwd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Cadillac'
       AND m.name = 'CTS-V'     -- model 'CTS-V'
   )
  ),
  -- Lexus LFA (2011) auction
  (43,
   'JTHDX5BH5B5001234',
   2011,
   5000,
   2,
   2,
   552,
   4805,
   'LEXLFA11',
   '2011-05-30',
   'white',
   'petrol',
   'rwd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Lexus'
       AND m.name = 'LFA'       -- model 'LFA'
   )
  ),
  -- Audi R8 (2014) auction
  (44,
   'WP0AB2984EL123456',
   2014,
   22000,
   2,
   2,
   430,
   4200,
   'AUDR814',
   '2014-07-10',
   'gray',
   'petrol',
   'awd',
   'automatic',
   7,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Audi'
       AND m.name = 'R8'        -- model 'R8'
   )
  ),
  -- Nissan GT-R (2015) auction
  (45,
   'JN1AR5EF7FM123456',
   2015,
   15000,
   2,
   4,
   545,
   3799,
   'NISGTR15',
   '2015-06-05',
   'black',
   'petrol',
   'awd',
   'dual_clutch',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Nissan'
       AND m.name = 'GT-R'      -- model 'GT-R'
   )
  ),
  -- McLaren 570S (2018) auction
  (46,
   'SBM13DAA8JW123456',
   2018,
   8000,
   2,
   2,
   562,
   3799,
   'MCL570S18',
   '2018-07-20',
   'orange',
   'petrol',
   'rwd',
   'automatic',
   7,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'McLaren Automotive'
       AND m.name = '570S Coupe'  -- model '570S Coupe'
   )
  ),
  -- Bugatti Veyron (2008) auction
  (47,
   'VF9SA2DAX18M12345',
   2008,
   6000,
   2,
   2,
   1001,
   7993,
   'BUGVY08',
   '2008-05-10',
   'blue',
   'petrol',
   'awd',
   'automatic',
   7,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Bugatti'
       AND m.name = 'Veyron'    -- model 'Veyron'
   )
  ),
  -- Porsche Cayman (2017) draft auction
  (48,
   'WP0AA2A97HS123456',
   2017,
   12000,
   2,
   2,
   345,
   2979,
   'PCAY2017',
   '2017-06-12',
   'gray',
   'petrol',
   'rwd',
   'manual',
   6,
   (SELECT m.id
      FROM models m
      JOIN manufacturers mf ON m.manufacturer_id = mf.id
     WHERE mf.name = 'Porsche'
       AND m.name = 'Cayman'   -- model 'Cayman'
   )
  );


INSERT INTO cars (
  offer_id,
  vin,
  production_year,
  mileage,
  number_of_doors,
  number_of_seats,
  engine_power,
  engine_capacity,
  registration_number,
  registration_date,
  color,
  fuel_type,
  drive,
  transmission,
  number_of_gears,
  model_id
) VALUES
  (1,  'JTDBR32E68A123456', 2018, 30000, 4, 5, 132, 1798, 'COR2018',    '2018-03-15', 'white',  'petrol',  'fwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Toyota'     AND m.name = 'Corolla')),
  (2,  '2HGFC2F59KH123456', 2019, 45000, 4, 5, 158, 1998, 'CIV2019',    '2019-07-10', 'black',  'petrol',  'fwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Honda'      AND m.name = 'Civic')),
  (3,  'WAUZZZF41LA123456', 2020, 20000, 4, 5, 190, 1984, 'AUDA42020',  '2020-05-05', 'gray',   'petrol',  'awd', 'automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Audi'       AND m.name = 'A4')),
  (4,  'WBA3A9G52FNS12345', 2017, 50000, 4, 5, 180, 1998, 'BMW3SP17',   '2017-09-12', 'blue',   'petrol',  'rwd', 'automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'BMW'        AND m.name = '3 Series')),
  (5,  'WVWZZZ1KZHW123456', 2019, 30000, 4, 5, 228, 1984, 'VWGTI2019',  '2019-04-20', 'red',    'petrol',  'fwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Volkswagen' AND m.name = 'GTI')),
  (6,  'WDDGF8AB8FA123456', 2016, 60000, 4, 5, 241, 1991, 'MERC16',     '2016-11-30', 'black',  'petrol',  'rwd', 'automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Mercedes-Benz' AND m.name = 'C-Class')),
  (7,  'JMZND2HT6M0123456', 2021, 15000, 2, 2, 155, 1998, 'MAZMX520',   '2021-02-14', 'red',    'petrol',  'rwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Mazda'      AND m.name = 'MX-5')),
  (8,  '1FA6P8TH0J5101234', 2018, 40000, 2, 4, 450, 5000, 'MUST2018',   '2018-08-22', 'red',    'petrol',  'rwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Ford'       AND m.name = 'Mustang')),
  (9,  'KM8J3CA46KU123456', 2020, 10000, 4, 5, 175, 1999, 'HYUTU20',    '2020-04-18', 'white',  'petrol',  'awd', 'automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Hyundai'    AND m.name = 'Tucson')),
  (10, 'KNDPMCAC7K7134567', 2019, 60000, 4, 5, 136, 1598, 'KIASPT19',   '2019-06-21', 'green',  'diesel',  'fwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Kia'        AND m.name = 'Sportage')),
  (11, 'VF1KZ0BA5HJ123456', 2017, 70000, 4, 5, 115, 1461, 'RNMNE17',    '2017-03-25', 'beige',  'diesel',  'fwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Renault'    AND m.name = 'Megane')),
  (12, 'W0L0XCF6861234567', 2020, 20000, 4, 5, 150, 1398, 'OPASTR20',   '2020-01-11', 'white',  'petrol',  'fwd', 'manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Opel'       AND m.name = 'Astra')),
  (13, '5YJ3E1EB8MF123456', 2021, 10000, 4, 5, 283,    1, 'TSLM32021',  '2021-01-01', 'white',  'electric','rwd','automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Tesla'      AND m.name = 'Model 3')),
  (14, '1N4AZ1CP4KC123456', 2019, 30000, 4, 5, 147,    1, 'NSILEF19',   '2019-03-30', 'blue',   'electric','fwd','automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Nissan'     AND m.name = 'Leaf')),
  (15, 'TMBJF9XK1K0123456', 2018, 60000, 4, 5, 115, 1968, 'SKOCT18',    '2018-06-12', 'gray',   'diesel',  'fwd','manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Škoda'      AND m.name = 'Octavia')),
  (16, 'VF3YHKEUAJN123456', 2020, 30000, 4, 5, 180, 1598, 'PEU300820',  '2020-09-15', 'green',  'petrol',  'fwd','automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Peugeot'    AND m.name = '3008')),
  (17, 'ZFA3120000B123456', 2019, 25000, 2, 4, 101, 1242, 'FT50019',    '2019-02-10', 'yellow', 'petrol',  'fwd','manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Fiat'       AND m.name = '500')),
  (18, 'YV4A22RK6H1123456', 2017, 50000, 4, 5, 250, 1969, 'VOVX6017',   '2017-07-07', 'black',  'petrol',  'awd','automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Volvo'      AND m.name = 'XC60')),
  (19, '3VW217AU1HM123456', 2020, 25000, 4, 5, 150, 1395, 'SEALEO20',   '2020-05-20', 'red',    'petrol',  'fwd','manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Seat'       AND m.name = 'Leon')),
  (20, 'VF7U7MPB4JL123456', 2018, 40000, 4, 5, 130, 1598, 'CTC418',     '2018-10-30', 'beige',  'petrol',  'fwd','manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Citroen'    AND m.name = 'C4')),
  (21, 'JTMRWRFV4KJ123456', 2019, 35000, 4, 5, 219, 2487, 'TYRAV419',   '2019-08-05', 'green',  'hybrid',  'awd','automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Toyota'     AND m.name = 'RAV4 Hybrid AWD')),
  (22, '1FADP3K29JL123456', 2020, 30000, 4, 5, 160, 1996, 'FDOCS20',    '2020-03-01', 'orange', 'petrol',  'fwd','manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Ford'       AND m.name = 'Focus')),
  (23, '3C4NJDEB1JT123456', 2018, 45000, 4, 5, 177, 1995, 'JPCOMP18',   '2018-05-10', 'brown',  'petrol',  'awd','automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Jeep'       AND m.name = 'Compass')),
  (24, 'WMWXY7C36M2L12345', 2021, 12000, 2, 4, 134, 1499, 'MICOOP21',   '2021-04-18', 'purple', 'petrol',  'fwd','automatic', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'MINI'       AND m.name = 'Cooper')),
  (25, 'UUCDCKKXKKA123456', 2019, 50000, 4, 5, 115, 1598, 'DADUS19',    '2019-11-22', 'maroon', 'diesel',  'fwd','manual', 6,
       (SELECT m.id FROM models m JOIN manufacturers mf ON m.manufacturer_id = mf.id
         WHERE mf.name = 'Dacia'      AND m.name = 'Duster'));


INSERT INTO reviews (reviewer_id, reviewee_id, description, rating, review_date) VALUES
  (2,  3, 'Smooth transaction and excellent communication',                    5, '2025-01-12 10:15:00'),
  (2,  4, 'Vehicle as described, would recommend',                             4, '2025-01-20 16:05:00'),
  (3,  2, 'Quick payment, friendly buyer',                                    5, '2025-02-05 09:45:00'),
  (3,  5, 'Minor delays in responses, but overall good experience',           4, '2025-02-10 11:30:00'),
  (4,  6, 'Car was well maintained, seller was helpful',                      5, '2025-03-01 14:10:00'),
  (4,  7, 'Slightly late delivery, but condition matched listing',            4, '2025-03-05 13:50:00'),
  (5,  4, 'Professional seller, smooth pickup',                               5, '2025-03-15 12:20:00'),
  (5,  8, 'Excellent handling of paperwork, quick responses',                 5, '2025-03-20 15:00:00'),
  (6,  5, 'Great communication, but pickup was delayed',                      3, '2025-04-01 10:00:00'),
  (6,  9, 'Pleasant experience, vehicle in good shape',                       5, '2025-04-05 09:30:00'),
  (7,  6, 'Seller was honest about condition, thanks',                        5, '2025-04-15 14:45:00'),
  (7, 10, 'Responsive and friendly, no issues',                               5, '2025-04-20 11:15:00'),
  (8,  2, 'Buyer asked good questions, prompt payment',                       5, '2025-05-01 16:00:00'),
  (8,  7, 'Quick transaction, friendly buyer',                                5, '2025-05-05 10:20:00'),
  (9,  8, 'Smooth process, highly recommended',                               5, '2025-05-10 14:00:00'),
  (9,  3, 'Fast payment, great communication',                                5, '2025-05-15 12:30:00'),
  (10, 9, 'Very polite, quick payment',                                       5, '2025-06-01 09:00:00'),
  (10, 2, 'Transaction completed smoothly',                                   5, '2025-06-05 13:45:00');


INSERT INTO users (name, email, password, role)
VALUES
    ('John Doe', 'john@example.com', 'password123', 'client'),
    ('Jane Smith', 'jane@example.com', 'password123', 'client')
ON CONFLICT (email) DO NOTHING;

INSERT INTO destinations (name, description, image_url)
VALUES
    (
        'Dushanbe',
        'The capital of Tajikistan, known for leafy boulevards, museums, tea houses, and lively evening promenades.',
        'https://images.unsplash.com/photo-1516483638261-f4dbaf036963?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Khujand',
        'A historic northern city on the Syr Darya with bustling bazaars, fortifications, and strong Silk Road character.',
        'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Iskanderkul',
        'A high-altitude alpine lake framed by dramatic peaks, waterfalls, and some of the best roadside scenery in Tajikistan.',
        'https://images.unsplash.com/photo-1501785888041-af3ef285b470?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Seven Lakes',
        'A chain of vividly colored lakes in the Fann Mountains, perfect for scenic drives, hiking, and village stays.',
        'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Pamir Highway',
        'One of the world’s great road journeys, crossing immense high-altitude landscapes and remote mountain communities.',
        'https://images.unsplash.com/photo-1521295121783-8a321d551ad2?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Khorog',
        'The heart of GBAO, combining river-valley calm, mountain views, botanical gardens, and a strong local culture.',
        'https://images.unsplash.com/photo-1482192505345-5655af888cc4?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Wakhan Valley',
        'A striking corridor of fortresses, hot springs, and sweeping views toward the Afghan and Pamiri mountains.',
        'https://images.unsplash.com/photo-1464820453369-31d2c0b651af?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Fann Mountains',
        'Jagged summits, turquoise lakes, and rewarding trekking routes make this one of Tajikistan’s signature adventure regions.',
        'https://images.unsplash.com/photo-1463694775559-eea25626346b?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Hisor Valley',
        'Home to the famous Hisor Fortress, orchard landscapes, and easy day-trip access from Dushanbe.',
        'https://images.unsplash.com/photo-1518684079-3c830dcef090?auto=format&fit=crop&w=1200&q=80'
    ),
    (
        'Penjikent',
        'A gateway to ancient Sogdian history, archaeological sites, and access points into the western Fann region.',
        'https://images.unsplash.com/photo-1506744038136-46273834b3fb?auto=format&fit=crop&w=1200&q=80'
    )
ON CONFLICT (name) DO UPDATE
SET
    description = EXCLUDED.description,
    image_url = EXCLUDED.image_url;

WITH seed_tours AS (
    SELECT *
    FROM (
        VALUES
            (
                'Dushanbe',
                'Dushanbe City Highlights',
                'A polished introduction to the capital covering museums, Rudaki Park, tea houses, and evening city life.',
                320.00,
                '2026-05-05 09:00:00'::timestamp,
                '2026-05-07 18:00:00'::timestamp,
                20
            ),
            (
                'Dushanbe',
                'Dushanbe Culture & Cuisine',
                'A slower city break focused on markets, Tajik food, local craft shops, and the capital’s modern cultural spaces.',
                360.00,
                '2026-06-14 10:00:00'::timestamp,
                '2026-06-16 17:00:00'::timestamp,
                18
            ),
            (
                'Dushanbe',
                'Capital Weekend Escape',
                'A relaxed urban itinerary with boulevard walks, monuments, cafes, and optional day excursions nearby.',
                295.00,
                '2026-09-18 11:00:00'::timestamp,
                '2026-09-20 16:00:00'::timestamp,
                22
            ),
            (
                'Khujand',
                'Khujand Silk Road Weekend',
                'Explore Panchshanbe Bazaar, the old citadel, the riverfront, and the rich mercantile history of northern Tajikistan.',
                340.00,
                '2026-05-21 09:00:00'::timestamp,
                '2026-05-24 18:00:00'::timestamp,
                18
            ),
            (
                'Khujand',
                'Khujand Heritage Trail',
                'A cultural route through museums, historic squares, local cuisine, and stories from the ancient trade road.',
                375.00,
                '2026-07-09 09:00:00'::timestamp,
                '2026-07-12 17:00:00'::timestamp,
                16
            ),
            (
                'Khujand',
                'Northern Tajikistan Discovery',
                'A broader city-and-surroundings tour blending Khujand highlights with village stops and regional food experiences.',
                410.00,
                '2026-10-02 08:00:00'::timestamp,
                '2026-10-06 18:00:00'::timestamp,
                15
            ),
            (
                'Iskanderkul',
                'Iskanderkul Lake Escape',
                'A scenic alpine retreat featuring lakeside walks, the nearby waterfall, and classic Fann Mountain panoramas.',
                540.00,
                '2026-05-28 08:00:00'::timestamp,
                '2026-05-31 17:00:00'::timestamp,
                16
            ),
            (
                'Iskanderkul',
                'Iskanderkul Road Adventure',
                'A drive-focused journey through mountain passes and viewpoints ending with quiet time at the lake.',
                575.00,
                '2026-07-24 07:00:00'::timestamp,
                '2026-07-27 18:00:00'::timestamp,
                14
            ),
            (
                'Iskanderkul',
                'Alpine Relaxation at Iskanderkul',
                'A slower-paced mountain stay built around fresh air, short hikes, and high-altitude scenery.',
                520.00,
                '2026-09-04 09:00:00'::timestamp,
                '2026-09-07 16:00:00'::timestamp,
                18
            ),
            (
                'Seven Lakes',
                'Seven Lakes Scenic Circuit',
                'Travel through all seven lakes with time for photos, village meals, and peaceful mountain viewpoints.',
                560.00,
                '2026-06-03 08:00:00'::timestamp,
                '2026-06-06 17:00:00'::timestamp,
                16
            ),
            (
                'Seven Lakes',
                'Haft Kul Trek & Village Stay',
                'A gentle active itinerary combining short walks, local hospitality, and changing lake colors from one stop to the next.',
                615.00,
                '2026-08-07 08:00:00'::timestamp,
                '2026-08-11 18:00:00'::timestamp,
                14
            ),
            (
                'Seven Lakes',
                'Seven Lakes Photography Route',
                'Designed for travelers chasing reflections, dramatic ridgelines, and golden-hour landscapes across the valley.',
                590.00,
                '2026-10-09 09:00:00'::timestamp,
                '2026-10-12 17:00:00'::timestamp,
                12
            ),
            (
                'Pamir Highway',
                'Pamir Highway Expedition',
                'A flagship overland route through high passes, windswept plateaus, and legendary long-distance road scenery.',
                1850.00,
                '2026-06-15 07:00:00'::timestamp,
                '2026-06-24 20:00:00'::timestamp,
                10
            ),
            (
                'Pamir Highway',
                'Roof of the World Drive',
                'A cinematic multi-day journey for travelers who want the full scale of Tajikistan’s eastern mountain landscapes.',
                1995.00,
                '2026-07-20 07:00:00'::timestamp,
                '2026-07-29 20:00:00'::timestamp,
                10
            ),
            (
                'Pamir Highway',
                'Pamir Road Adventure Lite',
                'A shorter version of the classic road trip focused on the most dramatic segments and remote overnight stops.',
                1625.00,
                '2026-09-12 08:00:00'::timestamp,
                '2026-09-19 19:00:00'::timestamp,
                12
            ),
            (
                'Khorog',
                'Khorog Valley Retreat',
                'A calm Pamiri town itinerary built around the botanical garden, river views, and local hospitality.',
                760.00,
                '2026-05-17 09:00:00'::timestamp,
                '2026-05-21 17:00:00'::timestamp,
                14
            ),
            (
                'Khorog',
                'Khorog & Mountain Life',
                'A balanced route blending Khorog’s urban core with nearby villages, markets, and scenic day drives.',
                845.00,
                '2026-08-13 08:00:00'::timestamp,
                '2026-08-18 18:00:00'::timestamp,
                12
            ),
            (
                'Khorog',
                'Pamiri Culture Week',
                'A culture-first itinerary with regional food, traditional music encounters, and immersive local experiences.',
                890.00,
                '2026-10-15 09:00:00'::timestamp,
                '2026-10-20 17:00:00'::timestamp,
                12
            ),
            (
                'Wakhan Valley',
                'Wakhan Valley Fortresses',
                'Trace the valley’s famous ancient fortresses while following the Panj River through one of Tajikistan’s most dramatic corridors.',
                1180.00,
                '2026-06-10 08:00:00'::timestamp,
                '2026-06-15 18:00:00'::timestamp,
                12
            ),
            (
                'Wakhan Valley',
                'Wakhan Hot Springs & Villages',
                'A slower route focused on rural stays, natural hot springs, and constant mountain views across the borderlands.',
                1095.00,
                '2026-07-30 08:00:00'::timestamp,
                '2026-08-04 17:00:00'::timestamp,
                12
            ),
            (
                'Wakhan Valley',
                'Wakhan Cultural Passage',
                'A storytelling-rich journey through shrines, village life, and the layered history of the far eastern valleys.',
                1140.00,
                '2026-09-24 09:00:00'::timestamp,
                '2026-09-29 18:00:00'::timestamp,
                10
            ),
            (
                'Fann Mountains',
                'Fann Mountains Basecamp',
                'A mountain adventure from a comfortable base with hikes, lake views, and star-filled nights.',
                890.00,
                '2026-06-19 08:00:00'::timestamp,
                '2026-06-24 18:00:00'::timestamp,
                14
            ),
            (
                'Fann Mountains',
                'Fann Lakes Trek',
                'A more active route through high meadows and alpine lakes for travelers who want classic Tajikistan trekking scenery.',
                980.00,
                '2026-07-16 07:00:00'::timestamp,
                '2026-07-22 19:00:00'::timestamp,
                12
            ),
            (
                'Fann Mountains',
                'Fann Peaks Explorer',
                'A panoramic multi-day escape focused on ridgeline views, camp nights, and big mountain atmosphere.',
                1040.00,
                '2026-09-03 08:00:00'::timestamp,
                '2026-09-09 18:00:00'::timestamp,
                10
            ),
            (
                'Hisor Valley',
                'Hisor Fortress Day Escape',
                'A short cultural tour centered on the fortress complex, local history, and nearby countryside.',
                210.00,
                '2026-04-25 09:00:00'::timestamp,
                '2026-04-26 18:00:00'::timestamp,
                24
            ),
            (
                'Hisor Valley',
                'Hisor & Valley Heritage',
                'A slightly longer route combining fortress visits, orchard landscapes, and village lunch stops.',
                285.00,
                '2026-05-30 09:00:00'::timestamp,
                '2026-06-01 17:00:00'::timestamp,
                20
            ),
            (
                'Hisor Valley',
                'Fortress & Countryside Weekend',
                'An easygoing weekend itinerary for travelers wanting history, landscapes, and a quick reset near Dushanbe.',
                295.00,
                '2026-10-23 10:00:00'::timestamp,
                '2026-10-25 16:00:00'::timestamp,
                20
            ),
            (
                'Penjikent',
                'Ancient Penjikent Explorer',
                'Discover archaeological sites, old city stories, and the deep historical roots of western Tajikistan.',
                430.00,
                '2026-05-08 09:00:00'::timestamp,
                '2026-05-11 17:00:00'::timestamp,
                16
            ),
            (
                'Penjikent',
                'Penjikent & Sogdian Heritage',
                'A heritage-first route through ruins, museums, and surrounding cultural landscapes.',
                465.00,
                '2026-08-21 09:00:00'::timestamp,
                '2026-08-24 18:00:00'::timestamp,
                14
            ),
            (
                'Penjikent',
                'Western Tajikistan Discovery',
                'A blended city-and-mountain-edge itinerary connecting Penjikent history with nearby scenic viewpoints.',
                520.00,
                '2026-09-25 08:00:00'::timestamp,
                '2026-09-29 18:00:00'::timestamp,
                14
            )
    ) AS v(destination_name, name, description, price, start_date, end_date, capacity)
)
INSERT INTO tours (destination_id, name, description, price, start_date, end_date, capacity)
SELECT
    d.id,
    st.name,
    st.description,
    st.price,
    st.start_date,
    st.end_date,
    st.capacity
FROM seed_tours st
JOIN destinations d ON d.name = st.destination_name
WHERE NOT EXISTS (
    SELECT 1
    FROM tours t
    WHERE t.name = st.name AND t.start_date = st.start_date
);

INSERT INTO bookings (user_id, tour_id, status)
SELECT u.id, t.id, b.status
FROM (
    VALUES
        ('john@example.com', 'Dushanbe City Highlights', 'confirmed'),
        ('john@example.com', 'Pamir Highway Expedition', 'pending'),
        ('jane@example.com', 'Seven Lakes Scenic Circuit', 'confirmed'),
        ('jane@example.com', 'Khujand Silk Road Weekend', 'cancelled')
) AS b(user_email, tour_name, status)
JOIN users u ON u.email = b.user_email
JOIN tours t ON t.name = b.tour_name
WHERE NOT EXISTS (
    SELECT 1
    FROM bookings existing
    WHERE existing.user_id = u.id AND existing.tour_id = t.id
);

INSERT INTO reviews (user_id, tour_id, rating, comment)
SELECT u.id, t.id, r.rating, r.comment
FROM (
    VALUES
        ('john@example.com', 'Dushanbe City Highlights', 5, 'A great introduction to the capital with smooth pacing and friendly guides.'),
        ('jane@example.com', 'Seven Lakes Scenic Circuit', 5, 'Absolutely stunning scenery and a very memorable mountain route.'),
        ('john@example.com', 'Pamir Highway Expedition', 4, 'Long days on the road, but the landscapes were completely worth it.')
    ) AS r(user_email, tour_name, rating, comment)
JOIN users u ON u.email = r.user_email
JOIN tours t ON t.name = r.tour_name
WHERE NOT EXISTS (
    SELECT 1
    FROM reviews existing
    WHERE existing.user_id = u.id AND existing.tour_id = t.id
);

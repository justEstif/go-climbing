-- Restore original learn content
UPDATE learn_content SET
  content = '<p>Silent feet is a technique where you focus on placing your feet on holds as quietly as possible. The goal is to eliminate the "clunk" of shoes hitting the wall, forcing you to use precision and intention with every foot placement.</p><p>Practice by climbing routes you know well and concentrating entirely on foot placement. Look at the hold, place your foot deliberately, and keep the weight transfer smooth.</p>',
  video_url = 'https://www.youtube.com/watch?v=kFrUiGks5Dc'
WHERE category = 'Footwork' AND title = 'Silent Feet';

UPDATE learn_content SET
  content = '<p>Smearing is the technique of using friction between your climbing shoe and the wall surface when there are no defined footholds. The rubber of your shoe sole creates enough friction to stand on if your weight is positioned correctly.</p><p>Key points: keep your heel low, push your weight into the wall, and trust your shoes.</p>',
  video_url = 'https://www.youtube.com/watch?v=sGXuAFiXHFI'
WHERE category = 'Footwork' AND title = 'Smearing';

UPDATE learn_content SET
  content = '<p>Edging uses the edge of your climbing shoe — typically the inside or outside edge — to stand precisely on small footholds. Inside edging (big toe side) is most common and gives the most power. Outside edging (pinky toe side) is used for drop-knee moves and traverses.</p><p>Always keep your heel slightly raised and drive your weight through the edge of the shoe for maximum contact area on the hold.</p>',
  video_url = NULL
WHERE category = 'Footwork' AND title = 'Edging Technique';

UPDATE learn_content SET
  content = '<p>The drop-knee (or Egyptian) is a technique where you turn your knee inward and downward, rotating your hip toward the wall. This allows you to reach farther with the opposite hand while keeping your center of gravity over your feet.</p><p>It is especially useful on overhanging terrain and for reducing the distance between your body and the holds above.</p>',
  video_url = 'https://www.youtube.com/watch?v=k5PtFhT_3w0'
WHERE category = 'Body Positioning' AND title = 'Hip Rotation (Drop-Knee)';

UPDATE learn_content SET
  content = '<p>Flagging is a technique where you extend one leg out to the side or behind you to counterbalance your body weight. It prevents you from barn-dooring (swinging away from the wall) when only one hand and one foot are on the same side.</p><p>Inside flag: leg crosses behind the weight-bearing foot. Outside flag: leg extends to the outside. Both keep your hips in and your center of gravity over your points of contact.</p>',
  video_url = 'https://www.youtube.com/watch?v=8JJSjFKFfxI'
WHERE category = 'Body Positioning' AND title = 'Flagging';

UPDATE learn_content SET
  content = '<p>The 4x4 is a classic power-endurance training method. Choose four boulder problems slightly below your max (around 60–70% effort). Climb all four back-to-back with minimal rest, then rest 3–4 minutes. Repeat for four rounds total.</p><p>This trains your body to maintain peak power output while fatigued — a critical skill for longer routes and competition climbing.</p>',
  video_url = 'https://www.youtube.com/watch?v=H0sNPWZ7e2Y'
WHERE category = 'Training Methods' AND title = 'The 4x4 Workout';

UPDATE learn_content SET
  content = '<p>Repeaters build finger strength and endurance using a hangboard. The classic protocol: hang for 7 seconds, rest 3 seconds, repeat 6 times per set. Rest 2–3 minutes between sets. Use an edge size that allows you to complete all reps — typically 20mm for intermediate climbers.</p><p>Focus on keeping your fingers in a half-crimp or open-hand position and avoid fully locking the crimp to reduce injury risk.</p>',
  video_url = 'https://www.youtube.com/watch?v=oGMPsS4oeIs'
WHERE category = 'Training Methods' AND title = 'Hangboard Repeaters';

UPDATE learn_content SET
  content = '<p>Crimps are small ledge holds that you grip with your fingers bent at the first knuckle. There are two main grip positions: full crimp (thumb wrapped over index finger — more power, higher injury risk) and half crimp (thumb not wrapped — safer for tendons).</p><p>Open-hand grip is the safest and should be trained regularly, even though it feels weaker initially. Strong open-hand grip is a hallmark of advanced climbers.</p>',
  video_url = NULL
WHERE category = 'Hold Types' AND title = 'Crimps';

UPDATE learn_content SET
  content = '<p>Slopers are rounded holds with no defined edge. They rely entirely on friction and require you to keep your palm as flat as possible against the surface, fingers pointing downward. Your center of gravity must stay low and directly below the hold.</p><p>Slopers punish poor footwork and body positioning — if your hips are away from the wall, you will slip off. They are grip-strength intensive and are particularly common on slab and outdoor granite routes.</p>',
  video_url = NULL
WHERE category = 'Hold Types' AND title = 'Slopers';

UPDATE learn_content SET
  content = '<p>Route reading is the practice of studying a route from the ground before you climb it. Look at the holds, identify likely rest positions, think through sequences, and plan your feet before your hands.</p><p>Good route readers can visualize the entire climb before touching the wall. Practice by watching other climbers, projecting each move mentally, and comparing your plan to what actually happens on the wall.</p>',
  video_url = 'https://www.youtube.com/watch?v=hLGiHjHY7GI'
WHERE category = 'Mental Skills' AND title = 'Route Reading';

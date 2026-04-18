/**
 * Usage:
 *   HORSE_ID=<uuid> TOKEN=<jwt> node scripts/seed-care-logs.mjs
 *
 * Get TOKEN from browser dev tools:
 *   Application → Local Storage → localhost:4200 → token (or auth_token)
 *
 * Get HORSE_ID from the URL when viewing the horse's detail page.
 */

const HORSE_ID = '0041c290-e094-4fdf-8a48-3a4b89ddeeeb';
const TOKEN    = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzY2MjQwNDAsImlhdCI6MTc3NjUzNzY0MCwic3ViIjoiZTkzY2Q3ZmEtYWU4Ny00ZDJkLWIzNzAtOThhYzU0YjQyMmVlIiwidHlwZSI6ImFjY2VzcyJ9.F9cdF8PivR5niIAe9VxXry85KPo73uVbNPwtyxsVobw';
const API      = process.env.API_URL ?? 'http://localhost:8080/api';

if (!HORSE_ID || !TOKEN) {
  console.error('Error: HORSE_ID and TOKEN env vars are required.');
  process.exit(1);
}

const logs = [
  { date: '2023-05-24', category: 'vaccination', notes: 'Received shots and coggins from Dr. Clark' },
  { date: '2023-06-25', category: 'vet',         notes: 'Showed signs of Gas Colic. Morning of didn\'t poop in stall. Laid down and got up repeatedly, pulse was elevated. Got casted and given banamine orally. Came out of cast. Taken to indoor arena to walk and seemed better. No longer showing signs of colic. Didn\'t eat breakfast, light dinner and no treats. Lunged (walk & trot) later that morning for ~10 minutes. Next morning pooped and back to normal.' },
  { date: '2023-08-04', category: 'vet',         notes: 'Showed signs of lameness in right front. Had a small jump lesson the day before, possibly sore from that.' },
  { date: '2023-08-07', category: 'vet',         notes: 'Continued to show lameness, very obvious in the right hind. Treated for an abscess for 3 days — leg was swollen, sensitive when picking foot, hoof was warm. Abscess blew a few months later.' },
  { date: '2023-08-13', category: 'other',       notes: 'First full body massage.' },
  { date: '2023-08-14', category: 'dental',      notes: 'Teeth floated by Tom.' },
  { date: '2023-08-16', category: 'vet',         notes: 'Right hind leg still swollen around the fetlock. Cold hosed the leg.' },
  { date: '2023-11-08', category: 'vet',         notes: 'Osteopathic consultation with Dr. Renschler. Full body exam. Needs teeth done properly. Shows signs of hind end gut and liver discomfort. Starting Immubiome G-Tract supplement and Milk Thistle. Feet are high-low. Starting journey to go barefoot — using Farrier\'s Fix to harden feet. Very reactive in hind and near poll; pain likely stems from diagonal teeth pattern, hind end gut, and liver issues. Scheduled teeth appointment.' },
  { date: '2023-11-08', category: 'diet',        notes: 'Changed feed to ½ quart ration balancer & 1 quart oats. Will monitor weight.' },
  { date: '2023-11-14', category: 'diet',        notes: 'Changed supplements: Phycox, Biotin Plus, ImmuBiome Lean Muscle (new), Milk Thistle (new).' },
  { date: '2023-11-28', category: 'vet',         notes: 'Follow up with Dr. Renschler. Did some acupuncture points.' },
  { date: '2023-11-29', category: 'diet',        notes: 'Started ⅓ scoop ImmuBiome G-Tract for PM supplement.' },
  { date: '2023-12-06', category: 'diet',        notes: 'Increased ImmuBiome G-Tract to ½ scoop for PM supplement.' },
  { date: '2023-12-13', category: 'diet',        notes: 'Increased ImmuBiome G-Tract to 1 full scoop for PM supplement.' },
  { date: '2023-12-19', category: 'dental',      notes: 'Dental with Dr. Renschler.' },
  { date: '2024-04-27', category: 'diet',        notes: 'Switched from oats and ration balancer to K.I.S Trace and oats.' },
  { date: '2024-05-02', category: 'vaccination', notes: 'Shots and coggins by Dr. Clark.' },
  { date: '2024-05-11', category: 'deworming',   notes: 'Spring deworming (Quest).' },
  { date: '2024-06-03', category: 'diet',        notes: 'Started Silver Linings Mare Moods supplement.' },
  { date: '2025-01-28', category: 'diet',        notes: 'Stopped KIS. Started Mad Barn Omneity. Stopped all other supplements. Feed is now oats, Kalm N\' EZ grain, and Omneity only.' },
  { date: '2025-02-02', category: 'other',       notes: 'Started 24/7 turnout.' },
  { date: '2025-02-14', category: 'dental',      notes: 'Dental with Dr. Renschler.' },
  { date: '2025-02-17', category: 'diet',        notes: 'Started giving some soaked alfalfa in the mornings.' },
  { date: '2025-02-21', category: 'diet',        notes: 'Started adding supplements back in with 1 qt soaked alfalfa: Salt, MSM, Mare Moods.' },
  { date: '2025-02-23', category: 'diet',        notes: 'Started adding ¼ cup alfalfa pellets to dry feed.' },
  { date: '2025-03-02', category: 'diet',        notes: 'Started GastroElm (1 tbsp) with 4 oz aloe vera juice to target ulcers. AM/PM until less ulcery, then switch to 1 tbsp powder in PM until bag is finished. If she responds well great; if not, try Gut-X. Goal: treat ulcers before saddle fitting.' },
  { date: '2025-03-04', category: 'other',       notes: 'First ride since November 2024. Rode bareback. Short ride — 10 minutes of walking. A little cinchy still. For first ride in four months with spring energy and a different bit, she did great.' },
  { date: '2025-04-18', category: 'deworming',   notes: 'Dewormed with Zimecterin Gold (Ivermectin & Praziquantel Paste).' },
  { date: '2025-05-02', category: 'other',       notes: 'Moved to new barn.' },
  { date: '2025-05-20', category: 'other',       notes: 'Saddle fitting with Morgan Fields.' },
  { date: '2025-06-01', category: 'diet',        notes: 'Stopped GastroElm.' },
  { date: '2025-06-21', category: 'vet',         notes: 'Showed some lameness in right hind leg. Had gone on a trail ride two days before. Cold hosed and kept an eye on it.' },
  { date: '2025-06-22', category: 'other',       notes: 'Added front pads back on. Ordered hoof boots for hind feet for trail rides. Received PEMF treatment to address soreness from previous day.' },
  { date: '2025-07-01', category: 'vet',         notes: 'Annual exam with Dr. Perry from Horse & Hound.' },
  { date: '2025-08-03', category: 'other',       notes: 'Finally have a proper saddle fitted.' },
  { date: '2025-10-06', category: 'vet',         notes: 'Eyelid swollen. Horse & Hound emergency call. Prescribed banamine and eye ointment. Swelling eventually went down over about 5 days.' },
  { date: '2025-10-12', category: 'other',       notes: 'English saddle flocked and fitted.' },
  { date: '2025-10-25', category: 'deworming',   notes: 'Fall deworming with Quest Plus Gel.' },
];

async function main() {
  console.log(`Seeding ${logs.length} care log entries for horse ${HORSE_ID}...\n`);
  let ok = 0, fail = 0;

  for (const log of logs) {
    const res = await fetch(`${API}/horses/${HORSE_ID}/care-logs`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${TOKEN}`,
      },
      body: JSON.stringify(log),
    });

    if (res.ok) {
      console.log(`✓  ${log.date}  [${log.category}]`);
      ok++;
    } else {
      const body = await res.json().catch(() => ({}));
      console.error(`✗  ${log.date}  [${log.category}]  →  ${res.status} ${JSON.stringify(body)}`);
      fail++;
    }
  }

  console.log(`\nDone: ${ok} succeeded, ${fail} failed.`);
}

main().catch(console.error);

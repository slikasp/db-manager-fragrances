# db-manager-fragrances

Database manager for a global list of fragrances, used in my FragranceTrackGo project.

The database itself is hosted on Supabase.

Used a fragrantica.com dataset from 2025 to start (24k records), found here: https://www.kaggle.com/datasets/olgagmiufana1/fragrantica-com-fragrance-dataset

Made some modifications for ease of use and developing this manager to update the database, expecting to get it up to ~100k records

<ol> Current plan to get the required details:
<li> Download all available cards (their links are the same except for ID) - DONE
<li> Scan the QRs in the cards to get the links to fragrances - IN PROGRESS

I can crop out the QR code from the image, but so far packages that can read them don't work - try others

<li> Parse HTML of the fragrance links for the rest of the details
</ol>

# db-manager-fragrances

Database manager for a global list of fragrances, used in my FragranceTrackGo project. The database itself is hosted on Supabase.

Used a fragrantica.com dataset from 2025 to start (24k records), found here: https://www.kaggle.com/datasets/olgagmiufana1/fragrantica-com-fragrance-dataset

Made some modifications for ease of use and developing this manager to update the database.
Currently over 80k records with existing cards (8 GB in size). Planning to have them locally and update if needed (another option is to fetche them from the source)
QR codes in the cards are stretched in both directions for some reason, so I have to fix them before decoding.
Getting remaining details from HTML will be a challenge, but should be pissible based on my research done for the Python equivalent.

<ol> Current plan to get the required details:
<li> Download all available cards (their links are the same except for ID) - DONE
<li> Scan the QRs in the cards to get the links to fragrances - DONE
<li> Parse HTML of the fragrance links for the rest of the details - IN PROGRESS
</ol>

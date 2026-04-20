# db-manager-fragrances

Database manager for a global list of fragrances. The database itself is hosted on Supabase, but I make snapshots available on Kaggle.com.

This is a database managet that maintains an up to date list of all fragrances in the world with various user submitted metrics from fragrantica.com.
It is parsing fragrance images with QR codes and looking for specific information in the HTML data of the encoded link.

I recommend building and running it on a server if you want to maintain your own database, first run will take some time (Currently there are over 80k records with cards, 8 GB in size).
<ol> What it does:
<li> Checks for any new fragrance cards that are not in cards table and downloads them
<li> Parses QR code on the downloaded cards and add new items in fragrances table
<li> Parses fragrance HTML for details and updates fragrances table
<li> Goes through all IDs that didn't have a card and checks if that changed
<li> Checks if there were updates to existing fragrances and updates database if needed

Current status:

<ol> Current plan to get the required details:
<li> Download all available cards (their links are the same except for ID) - DONE
<li> Scan the QRs in the cards to get the links to fragrances - DONE
<li> Parse HTML of the fragrance links for the rest of the details - IN PROGRESS
</ol>

The original dataset from 2024 that I started with (24k records) can be found here: https://www.kaggle.com/datasets/olgagmiufana1/fragrantica-com-fragrance-dataset
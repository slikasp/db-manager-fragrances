# db-manager-fragrances

Database manager for a global list of fragrances.
The database itself is hosted on Supabase (free tier), but I plan on making snapshots available on Kaggle.com.

This is a database manager that maintains an up to date list of all fragrances in the world with various user submitted metrics from fragrantica.com.
It is parsing fragrance images with QR codes and looking for specific information in the HTML data of the encoded link.

I mainly wrote this for myself to have an up to date fragrance list, but if you want to run it yourself I recommend doing it on a server, first run will take some time (Currently there are over 80k records with cards, >8 GB in size).
<ol> What it does:
<li> Goes through all fragrance IDs that didn't have a card and checks if they exist now
<li> Checks for any new fragrance cards that are not in cards table and downloads them
<li> Parses QR code on the downloaded cards and add new items in fragrances table
<li> Parses fragrance HTML for details and updates newly found fragrances
<li> (not yet implemented) Checks if there were updates to existing fragrances and updates database if needed
<ol>
Current status:

<ol> Planned work items:
<li> Download all available cards (their links are the same except for ID) - DONE
<li> Scan the QRs in the cards to get the links to fragrances - DONE
<li> Parse HTML of the fragrance links for the rest of the details - DONE (need to get around request spam)
<li> Make a service function that would run constantly (currently it's a one time job that does everything)
</ol>

The original dataset from 2024 that I started with (24k records) can be found here: https://www.kaggle.com/datasets/olgagmiufana1/fragrantica-com-fragrance-dataset
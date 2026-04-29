# db-manager-fragrances

Database manager for a list of fragrances.
The database itself is hosted remotely, not planning to upload it anywhere since it's not my data.

This is a database manager that maintains an up to date list of fragrances with various user submitted metrics from fragrantica.com.
It is parsing fragrance images with QR codes and looking for specific information in the HTML data of the encoded link.

I mainly wrote this for myself to have an up to date fragrance list, but if you want to run it yourself I recommend creating a list of ID's of fragrances you want to download. For example, I only need a list of ones I have owned myself, otherwise you hit fragrantica's traffic limits and all images take quite a bit of space.
<ol> What it does:
<li> Goes through fragrance IDs that didn't have a card and checks if they exist now
<li> Checks for any new fragrance cards that are not in cards table and downloads them
<li> Parses QR code on the downloaded cards and adds new items in fragrances table
<li> Parses fragrance HTML for details and updates newly found fragrances
<li> (not yet implemented) Checks if there were updates to existing fragrances and updates database if needed
<ol>

The original dataset from 2024 that I started with (24k records) can be found here: https://www.kaggle.com/datasets/olgagmiufana1/fragrantica-com-fragrance-dataset
# db-manager-fragrances

Database manager for a global list of fragrances.
The database itself is hosted remotely, not planning to upload it anywhere since it's not my data, but if this is something that would be useful to you - let me know.

## Why?

This is a database manager that maintains an up to date list of fragrances with various user submitted metrics from the web.
It keeps an up to date fragrance datebase by: donwloading specific images, decoding a QR code in them to figure out the website, searching that website for very specific details about the fragrance.

## How to use

I mainly wrote this for myself to have an up to date fragrance list for another app, so it's not really user friendly, if you want the access to the database - let me know. 
But if you really want to run it yourself I recommend creating a list of ID's of fragrances you want to download. For example, I only needed a list of ones I have owned myself, downloading all images take quite a bit of space, the full list is >120k items.
Also, run it as a service with low number of web requests, otherwise you'll hit website limits and get blocked.

It requires a config file with PostgreSQL database link in the executable directory to function (you can find the schema below):
`
{
    "remote_db_url": "postgresql://<db_address>:<db_port>/postgres"
}
`

## What it does:

<ol>
<li> Goes through fragrance IDs that didn't have a card and checks if they exist now</li>
<li> Checks for any new fragrance cards that are not in cards table and downloads them</li>
<li> Parses QR code on the downloaded cards and adds new items in fragrances table</li>
<li> Parses fragrance HTML for details and updates newly found fragrances</li>
<li> Updates newly added or oldest existing fragrances and downloads their card again needed</li>
</ol>

The original dataset from 2024 that I started with (24k records) can be found here: https://www.kaggle.com/datasets/olgagmiufana1/fragrantica-com-fragrance-dataset

## Database schema

### Table `cards`

Columns

| Name | Type | Constraints |
|------|------|-------------|
| `fragrantica_id` | `int4` | Primary |
| `url` | `text` |  |
| `image` | `text` |  |
| `has_card` | `bool` |  |
| `updated` | `timestamp` |  |

### Table `fragrances`

Columns

| Name | Type | Constraints |
|------|------|-------------|
| `id` | `int8` | Primary |
| `url` | `text` |  Nullable |
| `name` | `text` |  Nullable |
| `brand` | `text` |  Nullable |
| `country` | `text` |  Nullable |
| `gender` | `text` |  Nullable |
| `rating_value` | `numeric` |  Nullable |
| `rating_count` | `int4` |  Nullable |
| `year` | `int4` |  Nullable |
| `top_notes` | `text` |  Nullable |
| `middle_notes` | `text` |  Nullable |
| `base_notes` | `text` |  Nullable |
| `perfumer1` | `text` |  Nullable |
| `perfumer2` | `text` |  Nullable |
| `accord1` | `text` |  Nullable |
| `accord2` | `text` |  Nullable |
| `accord3` | `text` |  Nullable |
| `accord4` | `text` |  Nullable |
| `accord5` | `text` |  Nullable |
| `fragrantica_id` | `int4` |  Unique |
| `updated` | `timestamp` |  |
| `accord6` | `text` |  Nullable |
| `accord7` | `text` |  Nullable |
| `accord8` | `text` |  Nullable |
| `accord9` | `text` |  Nullable |
| `accord10` | `text` |  Nullable |

### Table `perfumers`

Columns

| Name | Type | Constraints |
|------|------|-------------|
| `name` | `text` | Primary |
| `country` | `text` |  |


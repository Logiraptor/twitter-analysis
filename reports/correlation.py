import json
import census
import pandas
from scipy.stats.stats import pearsonr
from pymongo import MongoClient
from matplotlib.backends.backend_pdf import PdfPages
import numpy as np


keys = {
    "TwoPop":            "P0030008",
    "TotalPop":          "P0030001",
    "WhitePop":          "P0030002",
    "BlackPop":          "P0030003",
    "AmericanIndianPop": "P0030004",
    "AsianPop":          "P0030005",
    "OtherPop":          "P0030007",
    "NativeHawaiianPop": "P0030006",
    # "Male":              "B01001_002E",
    # "Male<5":            "B01001_003E",
    # "Male5-9":           "B01001_004E",
    # "Male10-14":         "B01001_005E",
    # "Male15-17":         "B01001_006E",
    # "Male18-19":         "B01001_007E",
    # "Male20":            "B01001_008E",
    # "Male21":            "B01001_009E",
    # "Male22-24":         "B01001_010E",
    # "Male25-29":         "B01001_011E",
    # "Male30-34":         "B01001_012E",
    # "Male35-39":         "B01001_013E",
    # "Male40-44":         "B01001_014E",
    # "Male45-49":         "B01001_015E",
    # "Male50-54":         "B01001_016E",
    # "Male55-59":         "B01001_017E",
    # "Male60-61":         "B01001_018E",
    # "Male62-64":         "B01001_019E",
    # "Male65-66":         "B01001_020E",
    # "Male67-69":         "B01001_021E",
    # "Male70-74":         "B01001_022E",
    # "Male75-79":         "B01001_023E",
    # "Male80-84":         "B01001_024E",
    # "Male>85":           "B01001_025E",
    # "Female":            "B01001_026E",
    # "Female<5":          "B01001_027E",
    # "Female5-9":         "B01001_028E",
    # "Female10-14":       "B01001_029E",
    # "Female15-17":       "B01001_030E",
    # "Female18-19":       "B01001_031E",
    # "Female20":          "B01001_032E",
    # "Female21":          "B01001_033E",
    # "Female22-24":       "B01001_034E",
    # "Female25-29":       "B01001_035E",
    # "Female30-34":       "B01001_036E",
    # "Female35-39":       "B01001_037E",
    # "Female40-44":       "B01001_038E",
    # "Female45-49":       "B01001_039E",
    # "Female50-54":       "B01001_040E",
    # "Female55-59":       "B01001_041E",
    # "Female60-61":       "B01001_042E",
    # "Female62-64":       "B01001_043E",
    # "Female65-66":       "B01001_044E",
    # "Female67-69":       "B01001_045E",
    # "Female70-74":       "B01001_046E",
    # "Female75-79":       "B01001_047E",
    # "Female80-84":       "B01001_048E",
    # "Female>85":         "B01001_049E",
}

def prn(t):
	print t

client = MongoClient("localhost:27017")

tweetsCollection = client.clean_tweets.tweets
countyCollection = client.clean_tweets.counties

def projectRelevant(data):
	return {
		'id': data['_id'],
		'lat': data['computedcoords']['coordinates'][1],
		'lng': data['computedcoords']['coordinates'][0],
		'text': data['text'],
		'fips_state': data['fipsstate'],
		'fips_county': data['fipscounty'],
		'fips_tract': data['fipstract'],
		'fips_block': data['fipsblock'],
		'tweets': 1,
	}

all_tweets = list(tweetsCollection.find())
all_tweets = map(projectRelevant, all_tweets)

all_tweets = pandas.DataFrame(all_tweets)

by_state = all_tweets.groupby('fips_state').agg({'tweets': np.sum})
by_county = all_tweets.groupby(['fips_state', 'fips_county']).agg({'tweets': np.sum})
by_tract = all_tweets.groupby(['fips_state', 'fips_county', 'fips_tract']).agg({'tweets': np.sum})
# by_block = all_tweets.groupby(['fips_state', 'fips_county', 'fips_tract', 'fips_block']).agg({'tweets': np.sum})


print '%d tweets' % len(all_tweets)
print '%d states' % len(by_state)
print '%d counties' % len(by_county)
# print '%d tracts' % len(by_tract)
# print '%d blocks' % len(by_block)

for index, row in by_tract.iterrows():
	state, county, tract = index
	print census.resolveData(state, county, tract, keys)
	break
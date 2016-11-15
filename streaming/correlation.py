import json
import census
import pandas
from scipy.stats.stats import pearsonr
from pymongo import MongoClient
from matplotlib.backends.backend_pdf import PdfPages
import matplotlib.pyplot as plt
from pandas.tools.plotting import scatter_matrix
import numpy as np

client = MongoClient("localhost:27017")

tweetsCollection = client.clean_tweets.tweets
countyCollection = client.clean_tweets.counties

all_tweets = list(tweetsCollection.find())


def save_plot(plot, name):
    fig = plot.get_figure()
    fig.savefig(name)

keys = {
    "TwoPop":            "B02001_008E",
    "TotalPop":          "B02001_001E",
    "WhitePop":          "B02001_002E",
    "BlackPop":          "B02001_003E",
    "AsianPop":          "B02001_005E",
    "OtherPop":          "B02001_007E",
    "NativeHawaiianPop": "B02001_006E",
    "AmericanIndianPop": "B02001_004E",
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

def getCountyPop(state, county):
    ID = '%s:%s' % (state, county)
    resp = list(countyCollection.find({'_id':ID}))
    if resp:
        return resp[0]
    try:
        data = census.resolveData(state, county, keys)
        data = {key: int(data[0][key]) for key in keys.keys()}
    except Exception, e:
        data = {key: 0 for key in keys.keys()}
    countyCollection.update({'_id':ID}, data, upsert=True)
    return data


counties = {}
for tweet in all_tweets:
    state, countyCode = tweet['fipsstate'], tweet['fipscounty']

    key = state+":"+countyCode
    if key in counties:
        county = counties[key]
    else:
        county = getCountyPop(state, countyCode)
        county['Tweets'] = 0
        county['BroTweets'] = 0
        county['BruhTweets'] = 0
        county['BrahTweets'] = 0

    county['Tweets'] += 1
    county['BroTweets'] += (1 if tweet['bro'] else 0)
    county['BruhTweets'] += (1 if tweet['bruh'] else 0)
    county['BrahTweets'] += (1 if tweet['brah'] else 0)

    counties[key] = county


data = counties.values()

frame = pandas.DataFrame(data)


highConfidence = frame[frame['TotalPop'] > 0]


print 'With outliers:'
print np.sum(highConfidence['Tweets']), 'Tweets'
print np.sum(highConfidence['BroTweets']), 'BroTweets'
print np.sum(highConfidence['BruhTweets']), 'BruhTweets'
print np.sum(highConfidence['BrahTweets']), 'BrahTweets'
print len(highConfidence), 'Counties'

correlationHeaders = ['BruhTweets', 'BrahTweets', 'BroTweets', 'Tweets', 'TotalPop']
in_keys = keys.keys()


correlation = highConfidence.corr('pearson')
correlation = correlation[correlationHeaders]
interestingCorrelations = correlation.ix[in_keys]
print interestingCorrelations.sort(columns='TotalPop')


outlier_col = (highConfidence['TotalPop'])
highConfidence = highConfidence[~(outlier_col > (outlier_col.mean() + outlier_col.std()*2))]
print 'Without outliers:'
print np.sum(highConfidence['Tweets']), 'Tweets'
print np.sum(highConfidence['BroTweets']), 'BroTweets'
print np.sum(highConfidence['BruhTweets']), 'BruhTweets'
print np.sum(highConfidence['BrahTweets']), 'BrahTweets'
print len(highConfidence), 'Counties'
correlation = highConfidence.corr('pearson')
correlation = correlation[correlationHeaders]
interestingCorrelations = correlation.ix[in_keys]
print interestingCorrelations.sort(columns='TotalPop')


# print highConfidence.describe()

# pp = PdfPages('scatter.pdf')
# plots = [(x,y) for x in ['BruhTweets', 'BrahTweets', 'BroTweets', 'TotalPop'] for y in keys.keys()]
# for x, y in plots:
#     p = highConfidence.plot(x=x,y=y,kind='scatter')
#     pp.savefig(p.get_figure())
# pp.close()


# correlation = [(x, pearsonr(highConfidence[x], highConfidence[y]), y) for x in ['BruhTweets', 'BrahTweets', 'BroTweets', 'Tweets', 'TotalPop'] for y in keys.keys()]

# for x in sorted(correlation):
#     print x

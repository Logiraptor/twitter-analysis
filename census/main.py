import requests
from pymongo import MongoClient

census_key = "834f30d25c3352951df2bd4a457e21a88f9e083f"

# lat/lng to fips
# http://data.fcc.gov/api/block/find?latitude=39.9936&longitude=-105.0892&showall=false&format=json
# 080130608005010
# state(2):county(3):tract(6):block
# 08      :013      :060800  :5010

# fips to data
# for=tract:060800&in=state:08+county:013


def getStateCounty(lat, lng):
    r = requests.get(
        "http://data.fcc.gov/api/block/find?latitude=%f&longitude=%f&showall=false&format=json" % (lat, lng))
    data = r.json()
    return (data['State']['FIPS'], data['County']['FIPS'][2:]) if data['status'] == 'OK' else None


def getData(state, county, *variables):
    url = ("http://api.census.gov/data/2014/acs1?get=%s&for=county:%s&in=state:%s&key=" +
           census_key) % (",".join(variables), county, state)
    return requests.get(url).json()


def resolveData(state, county, format):
    keys, values = zip(*format.items())
    data = getData(state, county, *values)
    return map(lambda x: dict(zip(keys, x)), data[1:])

# state, county = getStateCounty(30.216452099999998, -92.0599479)

# print resolveData(state, county, {
#     "two": "B02001_008E",
#     "total": "B02001_001E",
#     "white": "B02001_002E",
#     "black": "B02001_003E",
#     "asian": "B02001_005E",
#     "other": "B02001_007E",
#     "native hawaiian": "B02001_006E",
#     "american indian": "B02001_004E",
# })

client = MongoClient()

tweets = client.engl452.tweets

coords = tweets.aggregate([
    {"$match": {"computedcoords": {"$ne": None}, "place.countrycode": "US"}},
    {"$project": {"coords": "$computedcoords.coordinates"}},
])

for coord in coords:
    print coord['coords'][1], coord['coords'][0]
    state, county = getStateCounty(coord['coords'][1], coord['coords'][0])
    print state, county


# So it'll look like this:
# Generate a collection where we count occurences of each form per block
# Normalize those values by population
# Compute significance with different ethnic trends

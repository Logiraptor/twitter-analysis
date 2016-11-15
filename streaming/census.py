import requests

census_key = "834f30d25c3352951df2bd4a457e21a88f9e083f"

def getData(state, county, *variables):
    other_resp = [[],[]]
    if len(variables) > 50:
        other_resp = getData(state, county, *variables[50:])
        variables = variables[:50]
    url = ("http://api.census.gov/data/2014/acs1?get=%s&for=county:%s&in=state:%s&key=" +
           census_key) % (",".join(variables), county, state)
    resp = requests.get(url)
    result = resp.json()

    return [a+b for a, b in zip(result, other_resp)]

    # [
    # [a, b]
    # [1, 2]
    # ]

    # [
    # [c, d]
    # [3, 4]
    # ]

    # [
    # [a, b, c, d]
    # [1, 2, 3, 4]
    # ]




def resolveData(state, county, format):
    keys, values = zip(*format.items())
    data = getData(state, county, *values)
    return map(lambda x: dict(zip(keys, x)), data[1:])

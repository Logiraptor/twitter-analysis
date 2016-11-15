import itertools
import nltk
from nltk.tokenize import TweetTokenizer
from pymongo import MongoClient
import nltk
from nltk.collocations import BigramCollocationFinder, TrigramCollocationFinder
bigram_measures = nltk.collocations.BigramAssocMeasures()
trigram_measures = nltk.collocations.TrigramAssocMeasures()


client = MongoClient("localhost:27017")

tweets = client.clean_tweets.tweets

all_tweets = list(tweets.find())

def show(all_tweets, key):
    print len(all_tweets), 'tweets'

    words = get_words(all_tweets)
    bi_grams = list(nltk.bigrams(words))
    tri_grams = list(nltk.trigrams(words))


    dist = nltk.FreqDist(bi_grams)

    output = ((val, count) for val, count in dist.most_common() if key in val)

    output = itertools.islice(output, 10)

    for val, count in output:
        print '\t', count, '\t', ' '.join(val)

    dist = nltk.FreqDist(tri_grams)

    output = ((val, count) for val, count in dist.most_common() if key in val)

    output = itertools.islice(output, 10)

    for val, count in output:
        print '\t', count, '\t', ' '.join(val)

excluded = ':,!'.split(',')

def get_words(tweets):
    text = ' '.join((t['text'] for t in tweets))

    tknzr = TweetTokenizer(preserve_case=False, reduce_len=True, strip_handles=True)

    words = tknzr.tokenize(text)
    return filter(lambda x: x not in excluded, words)

def collocations(words, key):
    creature_filter = lambda *x: key not in x
    finder = TrigramCollocationFinder.from_words(words)
    # only bigrams that appear 3+ times
    finder.apply_freq_filter(3)
    # only bigrams that contain 'creature'
    finder.apply_ngram_filter(creature_filter)
    # return the 10 n-grams with the highest PMI
    best = finder.nbest(trigram_measures.likelihood_ratio, 10)
    for val in best:
        print '\t', ' '.join(val)

# print 'all'
# show(all_tweets, '')

print 'bro'
collocations(get_words(filter(lambda x: x['bro'], all_tweets)), 'bro')
print 'bruh'
collocations(get_words(filter(lambda x: x['bruh'], all_tweets)), 'bruh')
print 'brah'
collocations(get_words(filter(lambda x: x['brah'], all_tweets)), 'brah')

# print 'bro'
# show(filter(lambda x: x['bro'], all_tweets), 'bro')
# print 'bruh'
# show(filter(lambda x: x['bruh'], all_tweets), 'bruh')
# print 'brah'
# show(filter(lambda x: x['brah'], all_tweets), 'brah')

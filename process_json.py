import numpy
import matplotlib.pyplot as plt
import matplotlib.lines as mlines
import json

FILE = "./experiment_container_2022.json"

with open(FILE) as f:
    data = json.load(f)

def convert_to_ms(string_val):
    if string_val.endswith("ms"):
       # do nothing, is ms
       return float(string_val.strip("ms"))
    if string_val.endswith("µs"):
       return float(string_val.strip("µs")) / 1000
    if string_val.endswith("s"):
       return float(string_val.strip("s")) * 1000

converted_data = {}
percentiles = {}

for key in data.keys():
    converted_data[key] = [ convert_to_ms(x.strip("[]")) for x in data[key].split() ]


for key in converted_data.keys():
    # print("VALUES FOR ", key)
    # print("MEAN VALUE: ", numpy.mean(converted_data[key]))
    percentiles[key] = {
        "95th": numpy.percentile(converted_data[key], 95),
        "50th": numpy.percentile(converted_data[key], 50)
    }
    mean = numpy.mean(converted_data[key])
    slowest = numpy.max(converted_data[key])
    fastest = numpy.min(converted_data[key]) 
    # print("95th percentile: ", numpy.percentile(converted_data[key], 95))
    # print("50th percentile", numpy.percentile(converted_data[key], 50))
    # print("SLOWEST: ", numpy.max(converted_data[key]))
    # print("FASTEST: ", numpy.min(converted_data[key]))
    # print(" -------------- ")

    print ("| %s | %ss | %ss | %ss | | %ss | %ss |" % (key, fastest / 1000, mean / 1000 , slowest / 1000 , percentiles[key]["95th"] / 1000 , percentiles[key]["50th"] / 1000))

    plt.figure(figsize=(10, 10), dpi=120)
    ax = plt.axes()
    ax.scatter(range(1,len(converted_data[key])+1), converted_data[key], marker=".")
    ax.plot([1,25,50,len(converted_data[key])], [percentiles[key]["95th"]]*4, color="red")
    ax.plot([1,25,50,len(converted_data[key])], [percentiles[key]["50th"]]*4, color="green")
    plt.title(("%s (HCSSHIM 0.9) Standard_D32s_v3" % key))
    plt.ylabel("time: ms")
    filename = ("./%s_WS2022_d32s_v3.png" % key)
    plt.savefig(filename, dpi=120, pad_inches=1, format="png")
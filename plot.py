from __future__ import print_function
import argparse
import numpy as np
import matplotlib.pyplot as plt
import csv
import sys

def plot_graph(file1, file2, labels, x_label, y_label, title, div, **kwargs):
    labels1 = []
    data1 = []
    for line in file1:
        data1.append(np.array([]))
        labels1.append(int(line[0]) / div)
        for i in range(1, len(line)):
            data1[len(data1) - 1] = np.append(data1[len(data1) - 1], float(line[i]))

    labels2 = []
    data2 = []
    for line in file2:
        data2.append(np.array([]))
        labels2.append(int(line[0]) / div)
        for i in range(1, len(line)):
            data2[len(data2) - 1] = np.append(data2[len(data2) - 1], float(line[i]))

    stdev1 = np.array([])
    stdev2 = np.array([])

    for i in range(len(data1)):
        stdev1 = np.append(stdev1, np.std(data1[i]))
        data1[i] = np.mean(data1[i])
    for i in range(len(data2)):
        stdev2 = np.append(stdev2, np.std(data2[i]))
        data2[i] = np.mean(data2[i])

    # if kwargs["blake_first"]:
        # label1 = "BLAKE3ZCC"
        # label2 = "SHA256"
    # else:
        # label1 = "SHA256"
        # label2 = "BLAKE3ZCC"
    fig, ax = plt.subplots()
    # ax.plot(labels1, data1, label=label1)
    # ax.plot(labels2, data2, label=label2)
    ax.errorbar(labels1, data1, xerr=None, yerr=stdev1, fmt='.-b', label=labels[0])
    ax.errorbar(labels2, data2, xerr=None, yerr=stdev2, fmt='.-r', label=labels[1])
    ax.set_xlabel(x_label)
    ax.set_ylabel(y_label)
    ax.set_title(title)
    ax.legend()

    plt.show()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    # if len(sys.argv) < 5:
        # print("Usage: python plot.py <file> <x label> <y label> <title>", file = sys.stderr)
        # sys.exit(1)
    parser.add_argument("--files", nargs=2)
    parser.add_argument("--labels", nargs=2)
    parser.add_argument("--x_label")
    parser.add_argument("--y_label")
    parser.add_argument("--title")
    parser.add_argument("--div", default=1024, type=int)
    args = parser.parse_args()
    with open(args.files[0], mode='r') as csv_file1:
        with open(args.files[1], mode='r') as csv_file2:
            csv_reader1 = csv.reader(csv_file1, delimiter=',')
            csv_reader2 = csv.reader(csv_file2, delimiter=',')
            plot_graph(file1=csv_reader1, file2=csv_reader2, \
                    labels=args.labels, x_label=args.x_label,
                    y_label=args.y_label, title=args.title, div=args.div)
                    # blake_first="true" in args.files[0], div=args.div)

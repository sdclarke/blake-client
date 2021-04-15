from __future__ import print_function
import argparse
import numpy as np
import matplotlib.pyplot as plt
import csv
import sys

colours = ['b', 'r', 'g', 'k', 'c', 'm']

def plot_graph(files, labels, x_label, y_label, title, div,
        storage=False, **kwargs):
    x = []
    y = []
    for f in files:
        x.append([])
        y.append([])
        for line in f:
            y[len(x)-1].append(np.array([]))
            x[len(x)-1].append(int(line[0]) / div)
            for i in range(1, len(line)):
                if storage:
                    y[len(y)-1][len(y[len(y)-1]) - 1] = np.append(y[len(y)-1][len(y[len(y)-1]) - 1],
                            float(line[i]) / div)
                else:
                    y[len(y)-1][len(y[len(y)-1]) - 1] = np.append(y[len(y)-1][len(y[len(y)-1]) - 1], float(line[i]))

    stdevs = [] # np.array([])

    if not storage:
        for j in range(len(y)):
            stdevs.append(np.array([]))
            for i in range(len(y[j])):
                stdevs[j] = np.append(stdevs[j], np.std(y[j][i]))
                y[j][i] = np.mean(y[j][i])
    fig, ax = plt.subplots()
    for i in range(len(y)):
        if storage:
            y0 = []
            y1 = []
            for j in range(len(y[i])):
                y0.append(y[i][j][0])
                y1.append(y[i][j][1])

            ax.errorbar(x[i], y0, xerr=None, yerr=None, fmt=f'P{colours[i]}', label=labels[i])
            ax.errorbar(x[i], y1, xerr=None, yerr=None, fmt=f'P{colours[i]}')
            ax.errorbar(x[i], list(map(lambda x,y: x+y, y0, y1)),
                    xerr=None, yerr=None, fmt=f'o--{colours[i+3]}',
                    label=f'Sum {labels[i]}', ms=4)
        else:
            ax.errorbar(x[i], y[i], xerr=None, yerr=stdevs[i], fmt=f'.-{colours[i]}', label=labels[i])
    ax.set_xlabel(x_label)
    ax.set_ylabel(y_label)
    ax.set_title(title)
    ax.legend()

    plt.show()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--files", nargs='+')
    parser.add_argument("--labels", nargs='+')
    parser.add_argument("--x_label")
    parser.add_argument("--y_label")
    parser.add_argument("--title")
    parser.add_argument("--div", default=1024, type=int)
    parser.add_argument("--storage", default=False, type=bool)
    args = parser.parse_args()
    files = []
    for f in args.files:
        csv_file = open(f, mode='r')
        files.append(csv.reader(csv_file, delimiter=','))
        
    plot_graph(files, labels=args.labels, x_label=args.x_label,\
        y_label=args.y_label, title=args.title, div=args.div,
        storage=args.storage)

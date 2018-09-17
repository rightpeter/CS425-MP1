#!/usr/bin/env python
# -*- coding: utf-8 -*-

import time
import os
import sys
import numpy as np
import matplotlib.pyplot as plt
import json

pattern_list = ["180", "74", "+", "[ib]", "[.]"]
grep_size = [35, 154, 592, 810, 932]
line_count = [204094, 1191425, 4440667, 6083248, 6999910]
repeat_time = 10


def grep_test(grep_pattern):
    """Test time consuming of single grep"""

    time_start = time.time()
    os.system("./client.bin {} > /dev/null".format(grep_pattern))
    time_end = time.time()
    return time_end - time_start


def benchmark():
    json_dic = {}
    json_dic['data'] = {}
    json_dic['time'] = []
    json_dic['err'] = []

    for pattern in pattern_list:
        print("Test grep: {}....\n".format(pattern))
        time_list = []
        for i in range(repeat_time):
            print("Test {}.....\n".format(i))
            time_list.append(grep_test(pattern))
        json_dic['data'][pattern] = time_list
        std_dev = np.std(time_list)
        avg = np.average(time_list)
        json_dic['time'].append(avg)
        json_dic['err'].append(std_dev)

    return json_dic


def main():
    print('sys: {}', sys.argv)
    json_dic = {}
    if '--benchmark' in sys.argv:
        json_dic = benchmark()

    if '--data' in sys.argv:
        with open('benchmark_data.json') as f:
            json_dic = json.load(f)

    if '--save_data' in sys.argv:
        with open('benchmark_data.json', 'w') as f:
            json.dump(json_dic, f)

    if '--plot' in sys.argv:
        plt.xlabel("Result Size (MB)")
        plt.ylabel("Time (s)")
        plt.grid()
        plt.errorbar(grep_size, json_dic['time'], json_dic['err'], ecolor='red', linestyle=':')
        for k in range(len(grep_size)):
            plt.annotate("({}, {:.2f})".format(grep_size[k], json_dic['time'][k]), (grep_size[k], json_dic['time'][k]))
        plt.show()


if __name__ == '__main__':
    main()

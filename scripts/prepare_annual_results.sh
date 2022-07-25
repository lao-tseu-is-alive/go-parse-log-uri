#!/bin/bash
time ./goParseLog data/big_one_5gb.log > data/log_gc/layers_visited_gc_space.log
cd data
sort layers_visited_gc_space.log > layers_visited_gc_space_sorted.log
uniq -c layers_visited_gc_space_sorted.log > layers_visited_gc_space_sorted_uniq_count.log
uniq layers_visited_gc_space_sorted.log > layers_visited_gc_space_sorted_uniq.log
gawk '$4~/2022/ {print $1}'  layers_visited_gc_space_sorted_uniq.log |uniq -c |sort -nr > layers_total_visits_2022.log
gawk '$4~/2021/ {print $1}'  layers_visited_gc_space_sorted_uniq.log |uniq -c |sort -nr > layers_total_visits_2021.log

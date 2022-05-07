# FEP methods on Ising models

Compute free energy difference between ferromagnet (J>0) and anti-ferromagnet (J<0).

$$
E = -J\sum_{<i,j>}S_i S_j - H \sum_i S_i
$$


## plan

1. sample 1 ensemble
    1. ferromagnetic and anti-ferromagnetic
    1. another choice of J, say J=0 and H=0
1. sample 2 ensembles, then process with BAR
1. sample more than 2 ensembles


## commands

- sample:
- exact:
- analyze: get averages and standard deviations from an ensemble

## TODO

- [  ] implement BAR

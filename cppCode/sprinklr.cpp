#include <bits/stdc++.h>
using namespace std;

#define ll long long

void solve(int V, vector<vector<pair<int, int>>> &graph, int K)
{
    priority_queue<pair<ll, int>, vector<pair<ll, int>>, greater<pair<ll,int>>> pq;
    vector<vector<ll>> dist(V + 1);

    pq.push({0, 1});

    while (!pq.empty())
    {
        int d = pq.top().first;
        int u = pq.top().second;
        pq.pop();

        if(dist[u].size()>=K) continue;
        dist[u].push_back(d);

        for(auto it:graph[u]){
            ll newDist = d+it.second;
            pq.push({newDist,it.first});
        }
    }
    for(int i = 0;i<K;i++){
        cout<<dist[V][i]<<" ";
    }
    return;
}

int main()
{
    freopen("input.txt", "r", stdin);   // Read from input.txt
    // freopen("output.txt", "w", stdout); // Optional: write output to output.txt

    int N, M, K;
    cin >> N >> M >> K;
    // Adjacency list: {destination, weight}
    vector<vector<pair<int, int>>> graph(N + 1);

    for (int i = 0; i < M; i++)
    {
        int u, v, w;
        cin >> u >> v >> w;
        graph[u].push_back({v, w});
    }

    // graph now contains all the edges
    solve(N, graph, K);
    return 0;
}
#include <bits/stdc++.h>
using namespace std;

void solve(int N, vector<int> &energy, vector<int> &loss)
{
    int ans = -1;

    for (int i = 0; i < N; i++)
    {
        if (energy[i] >= loss[i])
        {
            bool ok = true;
            long long sum = 0;

            for (int cnt = 0, j = i; cnt < N; cnt++, j = ((j + 1) % N))
            {
                sum += energy[j] - loss[j];
                if (sum < 0)
                {
                    ok = false;
                    break;
                }
            }

            if (ok)
            {
                cout << i << endl;
                return;
            }
        }
    }
    cout << ans << endl;
    return;
}

int main()
{
    freopen("input.txt", "r", stdin);
    // freopen("output.txt", "w", stdout);

    int N;
    cin >> N;

    vector<int> energy(N), loss(N);

    for (int i = 0; i < N; i++)
        cin >> energy[i];

    for (int i = 0; i < N; i++)
        cin >> loss[i];

    solve(N, energy, loss);

    return 0;
}
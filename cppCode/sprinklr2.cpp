#include <bits/stdc++.h>
using namespace std;

void solve(string s,string t){
    string s2,t2;
    
    for(char ch : s){
        if(ch != '#') s2+=ch;
        else if(s2.length()!=0){
            s2.pop_back();
        }
    }

    for(char ch : t){
        if(ch != '#') t2+=ch;
        else if(t2.length()!=0){
            t2.pop_back();
        }
    }

    cout<<((s2 == t2)?"true":"false")<<endl;

    return;

}

int main() {
    freopen("input.txt", "r", stdin);
    // freopen("output.txt", "w", stdout);

    string s, t;
    cin >> s >> t;

    // Your logic here

    solve(s,t);
    return 0;
}
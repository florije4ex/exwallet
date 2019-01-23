# Root key (for EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV)
# 5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3

# Root public key (EOS..5CV)
export owner_pubkey=EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
export active_pubkey=EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV

function cleos() {
  docker exec regtest_keosd cleos -u http://nodeosd:8888 --wallet-url http://localhost:8900 "$@"
}

function keosd() {
  docker exec regtest_nodeosd keosd "$@"
}

function pkill() {
  docker exec regtest_nodeosd pkill "$@"
}

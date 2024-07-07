var provider = window.ethereum;
var web3 = undefined;   

async function connectWallet() {
  try{
    let accounts = await web3.send("eth_accounts", [])
    if(accounts.length == 0){
      accounts = await web3.send("eth_requestAccounts", [])
      return
    }

    provider.emit("accountsChanged", accounts)
  }catch(err){  
    console.error("connectAccount", err)
    if(err.code === 4001){
      alert(`Failed to connect metamask account: ${err.message}`)
    }
  }
}

async function handleAccountsChanged(accounts) {
  try{
    console.log("handleAccountsChanged", accounts)
    const account = accounts.length > 0 ? accounts[0] : ""
    await htmx.ajax('POST', '/events/onAccountConnected', {
      target:'#navbar', 
      swap:'outerHTML', 
      values: { account }
    })
  }catch(err){
    console.error("handleAccountsChanged", err)
  }
}

function connectMetaMask() {
  const web3 = new ethers.providers.Web3Provider(provider)
  if(!provider.isConnected) {
    alert("Metamask is not connected to current chain, please reload the page and try again")
    return
  }

  provider.on("_initialized", (connectInfo) => console.warn("_initialized not implemented", connectInfo));
  provider.on("connect", (connectInfo) => console.warn("connect not implemented", connectInfo));
  provider.on("accountsChanged", handleAccountsChanged);
  provider.on("chainChanged", (chainId) => console.warn("chainChanged not implemented", chainId));
  provider.on("disconnect", (error) => console.warn("disconnect not implemented", error));

  return web3
}

function main(){
  if(provider=== undefined) {
    alert("No ethereum provider")
    return
  }
  if(typeof window.ethereum.isMetaMask === "boolean" && window.ethereum.isMetaMask){
    web3 = connectMetaMask(provider)
  }

  if(web3 === undefined){
    alert("No web3 provider")
    return;
  }
}

main()
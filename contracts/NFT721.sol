// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract NFT721 is ERC721URIStorage, Ownable {
	uint public tokenCount;

	mapping(uint => uint) public prices;
	mapping(uint => bool) public isListing;
	
	constructor(address initialOwner) ERC721("NFT", "NFT") Ownable(initialOwner){}

	function mint(string memory tokenURI, uint price) external onlyOwner returns (uint) {
		require(price > 0, "price must be a positive value");
		uint tokenId = tokenCount++;
		_safeMint(msg.sender, tokenId);
		_setTokenURI(tokenId, tokenURI);

		prices[tokenId] = price;
		isListing[tokenId] = true;
		return (tokenId);
	}

	function list(uint tokenId, uint price) external {
		require(msg.sender == ownerOf(tokenId), "Unauthorized");
		require(!isListing[tokenId], "Already listing");
		require(price > 0, "Price must be a positive value");
		prices[tokenCount] = price;
		isListing[tokenCount] = true;
	}

	function unlist(uint tokenId) external {
		require(msg.sender == ownerOf(tokenId), "Unauthorized");
		require(isListing[tokenId], "Not listing");
		isListing[tokenCount] = false;
	}

	function buy(uint tokenId) external payable {
		require(isListing[tokenId], "Not listing");
		require(msg.value == prices[tokenId], "Invalid payment value");

		isListing[tokenCount] = false;
		_safeTransfer(ownerOf(tokenId), msg.sender, tokenId);
	}
}

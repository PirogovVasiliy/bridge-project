// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract BridgeToken is ERC20, Ownable{
    constructor() ERC20("BridgeToken", "BT") Ownable(msg.sender){}

    event TokensSentToBridge(string to, uint amount);

    function sendTokensToBridge(string memory to, uint amount) external{
        require(amount > 0, "Amount of transaction can not be zero!");
        require (balanceOf(msg.sender) >= amount, "Not enough tokens!");
        _burn(msg.sender, amount);
        emit TokensSentToBridge(to, amount);
    }

    function receiveFromBridge(address to, uint amount)external onlyOwner{
        _mint(to, amount);
    }

    function sellTokens(uint amount)external{
        require(amount <= balanceOf(msg.sender), "Not enough tokens!");
        payable(msg.sender).transfer(amount);
        _burn(msg.sender, amount);
    }

    function buyTokens()external payable{
        assert(msg.value > 0);
        _mint(msg.sender, msg.value);
    }
}
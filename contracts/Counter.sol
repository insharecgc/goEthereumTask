// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract Counter {
    uint256 public count;

    event AddCount(uint256 newCount);

    constructor() {
        count = 1;
    }

    function addOne() external {
        count++;
        emit AddCount(count);
    }

}
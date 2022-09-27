// SPDX-License-Identifier: Apache-2.0

pragma solidity >=0.7.1;

contract CallerTest {

    event ReceiveRequest(bool success, bytes returndata);

    function call(address target, bytes memory data, uint256 value) public returns (bytes memory) {
        (bool success, bytes memory returndata) = target.call{value: value}(data);
        emit ReceiveRequest(success, returndata);
    }
}

pragma solidity >=0.7.1;


contract Bridge {

    event ReceiveRequestWithReturnValue(bool success, bytes sentData, string returndata);

    function receiveRequestWithReturnValue(bytes memory _sel, address receiveSide) external {
        string memory errorMessage;

        (bool success, bytes memory returndata) = receiveSide.call(_sel);
        if (!success) {
            errorMessage = sliceDestructive(returndata);
        }

        emit ReceiveRequestWithReturnValue(success, _sel, errorMessage);

    }


    function receiveRequestWithReturnValue2(bytes memory _sel, address receiveSide) external {
        string memory _return_data;

        (bool success, bytes memory returndata) = receiveSide.call(_sel);
        if (!success) {
            _return_data = sliceDestructive(returndata);
        }

        emit ReceiveRequestWithReturnValue(success, _sel, _return_data);

    }


    function receiveRequestAllwaysSuccess(bytes memory _sel, address receiveSide) external {
        emit ReceiveRequestWithReturnValue(true, _sel, string(_sel));
    }


    function receiveRequestWithVerifyCallResult(bytes memory _sel, address receiveSide) external {
        (bool success, bytes memory returndata) = receiveSide.call(_sel);
        emit ReceiveRequestWithReturnValue(success, _sel, string(returndata));
    }


    function delegateAndReturn(bytes memory _sel, address implementation) internal returns (bool, bytes memory) {
        (bool success, bytes memory ret) = implementation.call(_sel);
        if (!success) {
            assembly {

                let m := mload(0x40)
                returndatacopy(m, 0, returndatasize())
                return (0, returndatasize())
            }

        }
        return (success, ret);
    }

    function sliceDestructive(
        bytes memory input
    )
    internal
    pure
    returns (string memory result)
    {
        if (input.length > 0) {
            // Create a new bytes structure around [from, to) in-place.
            assembly {
                let len := returndatasize()
                result := add(input, 68)
                mstore(result, len)
            }
        }
        return result;
    }
}

pragma solidity >=0.7.1;


contract Bridge {

    event ReceiveRequestWithReturnValue(bool success, bytes sentData, bytes returndata);

    function receiveRequestWithReturnValue(bytes memory _sel, address receiveSide) external {
        bytes memory _return_data;

        (bool success, bytes memory returndata) = receiveSide.call(_sel);
        if (!success) {
            _return_data = sliceDestructive(returndata);
        }

        emit ReceiveRequestWithReturnValue(success, _sel, _return_data);

    }


    function receiveRequestWithReturnValue2(bytes memory _sel, address receiveSide) external {
        bytes memory _return_data;

        (bool success, bytes memory returndata) = receiveSide.call(_sel);
        if (!success) {
            assembly {
                let from := 68
                let len := returndatasize()
                let to := add(from, len)
                _return_data := add(returndata, from)
                mstore(_return_data, sub(to, from))
            }

            emit ReceiveRequestWithReturnValue(success, _sel, _return_data);

        }
    }



        function receiveRequestAllwaysSuccess(bytes memory _sel, address receiveSide) external {
            emit ReceiveRequestWithReturnValue(true, _sel, _sel);
        }


        function receiveRequestWithVerifyCallResult(bytes memory _sel, address receiveSide) external {
            (bool success, bytes memory returndata) = receiveSide.call(_sel);
            emit ReceiveRequestWithReturnValue(success, _sel, returndata);
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
    returns (bytes memory result)
    {
        // Create a new bytes structure around [from, to) in-place.
        assembly {
            let from := 68
            let len := returndatasize()
            let to := add(from, len)
            result := add(input, from)
            mstore(result, sub(to, from))
        }
        return result;
    }
}

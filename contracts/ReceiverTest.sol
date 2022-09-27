// SPDX-License-Identifier: Apache-2.0

pragma solidity >=0.7.1;

contract ReceiverTest {

    uint256 public testValue = 0;


    function setTestValue(uint256 _value) public {
        testValue = _value;
    }


    function failToSetTestValue(uint256 _value) public {
        require(false, "<Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus vel felis tortor. Mauris at scelerisque justo. Etiam vehicula justo ut hendrerit accumsan. Etiam vitae mauris ac purus ultrices scelerisque. Aliquam erat volutpat. Integer arcu urna, placerat vel orci vitae, vestibulum congue nisl. Sed congue laoreet volutpat. Vestibulum ut dictum ex, et viverra nisl. Quisque varius lacinia convallis. Aliquam eget leo at neque aliquet pretium a quis est. Sed volutpat risus justo, vel efficitur purus lacinia fringilla. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Sed volutpat cursus nulla, eu egestas mauris consectetur vitae.>");
    }


}


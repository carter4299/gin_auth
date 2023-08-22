import React, { useState } from 'react';
import { base64EncodeToBytes, base64Encode, validateUsername, validatePassword } from './auth_util';
import { authenticateUser } from './auth_api';

export function useLoginHooks(setRecheckAuth: any) {

    const [userName, setUserName] = useState('');
    const [userPass, setUserPass] = useState('');


    const handleUserNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        console.log("Returning User Name Recieved");
        console.log(e.target.value);
        setUserName(e.target.value);
    };
    const handleUserPassChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        console.log("Returning User Password Recieved");
        console.log(e.target.value);
        setUserPass(e.target.value);
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if(!validateUsername(userName) || !validatePassword(userPass)) {
            console.error('Failed validation');
            return;
        }

        const encodedUserName = base64Encode(userName);
        const encodedUserPass = base64EncodeToBytes(userPass);

        const response = await authenticateUser(encodedUserName, encodedUserPass);
        if (response.ok) {
            console.log('Successfully registered user');
            setRecheckAuth((prev: boolean) => !prev);
        } else {
            console.error('Failed to register user');
        }
    };


    return {
        userName, setUserName,
        userPass, setUserPass,
        handleUserNameChange,
        handleUserPassChange,
        handleSubmit,
    }

}



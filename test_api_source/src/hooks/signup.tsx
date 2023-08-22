import React, { useState } from 'react';
import { base64EncodeToBytes, base64Encode, validateUsername, validatePassword } from './auth_util';
import { registerUser } from './auth_api';

export function useSignupHooks(setRecheckAuth: any) {
    const [newUserName, setNewUserName] = useState('');
    const [newUserPass, setNewUserPass] = useState('');


    const handleNewUserNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        console.log("New User Name Received");
        console.log(e.target.value);
        setNewUserName(e.target.value);
    };
    const handleNewUserPassChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        console.log("New User Password Received");
        console.log(e.target.value);
        setNewUserPass(e.target.value);
    };

    const handleRegistration = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if(!validateUsername(newUserName) || !validatePassword(newUserPass)) {
            console.error('Failed validation');
            return;
        }

        const encodedNewUserName = base64Encode(newUserName);
        const encodedNewUserPass = base64EncodeToBytes(newUserPass);

        const response = await registerUser(encodedNewUserName, encodedNewUserPass);

        if (response.ok) {
            console.log('Successfully registered user');
            setRecheckAuth((prev: boolean) => !prev);
        } else {
            console.error('Failed to register user');
        }
    };

    return {
        newUserName, 
        newUserPass,
        handleNewUserNameChange,
        handleNewUserPassChange,
        handleRegistration
    }
}

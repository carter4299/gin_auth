import { createContext, useContext } from 'react';
type UserContextType = {
    recheckAuth: boolean;
    setRecheckAuth: React.Dispatch<React.SetStateAction<boolean>>;
};


export const UserContext = createContext<UserContextType | undefined>(undefined);

export const useUserContext = () => {
    const context = useContext(UserContext);
    if (!context) {
        throw new Error("useUserContext must be used within a UserContext.Provider");
    }
    return context;
};
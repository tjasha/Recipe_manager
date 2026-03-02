import React from 'react';

type MainLayoutProps = {
    children: React.ReactNode;
};

export default function MainLayout({ children }: MainLayoutProps) {
    return (
        <div>
            {children}
        </div>
    );
}

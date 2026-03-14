import React from 'react';
import NavigationBar from '../components/NavigationBar';

type MainLayoutProps = {
    children: React.ReactNode;
};

export default function MainLayout({ children }: MainLayoutProps) {
    return (
        <div>
            <NavigationBar />
            <main className="container">
                {children}
            </main>
        </div>
    );
}

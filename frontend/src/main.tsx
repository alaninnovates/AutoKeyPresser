import React from 'react'
import {createRoot} from 'react-dom/client'
import App from './App'
import {createTheme, MantineProvider} from '@mantine/core';

const theme = createTheme({

});
const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <MantineProvider theme={theme} forceColorScheme={'dark'}>
            <App/>
        </MantineProvider>
    </React.StrictMode>
)

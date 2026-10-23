import '@mantine/core/styles.css';
import {
    ActionIcon,
    AppShell,
    AppShellFooter,
    AppShellHeader,
    AppShellMain,
    Button,
    Kbd,
    Tabs,
    TabsList,
    TabsPanel,
    TabsTab,
    Text
} from "@mantine/core";
import {GearIcon} from '@phosphor-icons/react';
import {useState} from "react";
import {KeyPresserPanel} from "./KeyPresserPanel";
import {SetActiveTab} from "../bindings/AutoKeyPresser/appservice";
import {ActiveTab} from "../bindings/AutoKeyPresser";

function App() {
    const [activeTab, setActiveTab] = useState<string | null>('key_presser');

    return (
        <AppShell>
            <Tabs value={activeTab} onChange={async (value) => {
                const newAT = await SetActiveTab(value as any);
                setActiveTab(newAT);
            }}>
                <AppShellHeader p='xs' display="flex" style={{gap: "1rem", justifyContent: "space-between"}}>
                    <TabsList>
                        <TabsTab value={ActiveTab.ActiveTabKeyPresser}>Key Presser</TabsTab>
                        <TabsTab value={ActiveTab.ActiveTabSlime}>Slime Killer</TabsTab>
                        <TabsTab value={ActiveTab.ActiveTabFisher}>Auto Fisher</TabsTab>
                    </TabsList>
                    <ActionIcon size="lg" radius="xl"
                                onClick={() => console.log("Settings clicked")}>
                        <GearIcon/>
                    </ActionIcon>
                </AppShellHeader>
                <AppShellMain
                    pl='md' pr='md'
                    h="calc(100vh - var(--app-shell-header-height, 0px))"
                    display="flex"
                >
                    <TabsPanel value="key_presser" h="100%" display="flex" style={{flexDirection: 'column'}}>
                        <KeyPresserPanel />
                    </TabsPanel>
                    <TabsPanel value="slime" h="100%" display="flex" style={{flexDirection: 'column'}}>
                        <Text flex={1}>Slime Killer</Text>
                    </TabsPanel>
                    <TabsPanel value="fisher" h="100%" display="flex" style={{flexDirection: 'column'}}>
                        <Text>Auto Fisher</Text>
                    </TabsPanel>
                </AppShellMain>
            </Tabs>
            <AppShellFooter p='sm' display="flex" style={{gap: "1rem", justifyContent: "flex-end"}}>
                <Button>Start <Kbd ml={'sm'}>F1</Kbd></Button>
                <Button>Stop <Kbd ml={'sm'}>F3</Kbd></Button>
            </AppShellFooter>
        </AppShell>
    )
}

export default App

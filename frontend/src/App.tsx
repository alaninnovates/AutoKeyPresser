import '@mantine/core/styles.css';
import {
    ActionIcon,
    AppShell,
    AppShellFooter,
    AppShellHeader,
    AppShellMain,
    Button,
    Card,
    Group,
    Kbd,
    NumberInput,
    Radio,
    ScrollArea,
    Scroller,
    Stack,
    Switch,
    Tabs,
    TabsList,
    TabsPanel,
    TabsTab,
    Text
} from "@mantine/core";
import {GearIcon, KeyboardIcon, MouseIcon} from '@phosphor-icons/react';
import {useState} from "react";

const characters = [
    'Shift', 'Ctrl', 'Alt', 'Meta',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
    'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
    '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
    '`', '-', '=', '[', '\\', ';', '\'', ',', '.', '/', '~',
    '_', '+', '{', '|', ':', '"', '<', '>', '?'
];

interface MouseAction {
    type: 'mouse';
    button: 'left' | 'right' | 'middle';
    hold: boolean;
    repeat: boolean;
    repeatInterval: number;
}

interface KeyboardAction {
    type: 'keyboard';
    key: typeof characters[number];
    hold: boolean;
    repeat: boolean;
    repeatInterval: number;
}

type Action = MouseAction | KeyboardAction;

function App() {
    const [activeTab, setActiveTab] = useState<string | null>('key_presser');
    const [actions, setActions] = useState([
        {
            type: 'keyboard',
            key: 'Shift',
            hold: true,
            repeat: false,
            repeatInterval: 0,
        },
        {
            type: 'keyboard',
            key: 'Ctrl',
            hold: true,
            repeat: false,
            repeatInterval: 0,
        },
        {
            type: 'keyboard',
            key: 'd',
            hold: true,
            repeat: false,
            repeatInterval: 0,
        },
        {
            type: 'mouse',
            button: 'left',
            hold: false,
            repeat: true,
            repeatInterval: 1000,
        }
    ] as Action[]);
    return (
        <AppShell>
            <Tabs value={activeTab} onChange={setActiveTab}>
                <AppShellHeader p='xs' display="flex" style={{gap: "1rem", justifyContent: "space-between"}}>
                    <TabsList>
                        <TabsTab value="key_presser">Key Presser</TabsTab>
                        <TabsTab value="slime">Slime Killer</TabsTab>
                        <TabsTab value="fisher">Auto Fisher</TabsTab>
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
                        <Scroller flex={1}>
                            <Group wrap="nowrap" h="100%" align="stretch" gap="md">
                                {actions.map((action, index) => (
                                    <Card
                                        key={index}
                                        withBorder
                                        miw={200}
                                        h="100%"
                                        display="flex"
                                        style={{flexDirection: 'column'}}
                                    >
                                        <Text fw={500}
                                              mb="xs">{action.type === 'keyboard' ? 'Keyboard' : 'Mouse'} Action {index + 1}</Text>

                                        {action.type === 'keyboard' ? (
                                            <ScrollArea
                                                type="always"
                                                scrollbarSize={6}
                                                offsetScrollbars="y"
                                                h="150px"
                                            >
                                                <Stack>
                                                    {characters.map((char, charIndex) => (
                                                        <Text key={charIndex} size="sm">{char}</Text>
                                                    ))}
                                                </Stack>
                                            </ScrollArea>
                                        ) : (<>
                                            <Radio.Group
                                                label="Mouse Button"
                                                value={action.button}
                                                onChange={(value) => console.log(`Mouse button changed to ${value}`)}
                                                pb="sm"
                                            >
                                                <Stack mt="xs" gap='xs'>
                                                    <Radio value="left" label="Left"/>
                                                    <Radio value="right" label="Right"/>
                                                    <Radio value="middle" label="Middle"/>
                                                </Stack>
                                            </Radio.Group>
                                        </>)}
                                        <Text fw={400} mb="xs" size='sm'>Modifiers</Text>
                                        <Switch
                                            checked={action.hold}
                                            onChange={(event) => {
                                                setActions(prevActions => {
                                                    const newActions = [...prevActions];
                                                    newActions[index] = {
                                                        ...newActions[index],
                                                        hold: event.target.checked,
                                                        repeat: event.target.checked ? false : newActions[index].repeat,
                                                    };
                                                    return newActions;
                                                })
                                            }}
                                            label="Hold"
                                        />
                                        <Switch
                                            pt='md'
                                            checked={action.repeat}
                                            onChange={(event) => {
                                                setActions(prevActions => {
                                                    const newActions = [...prevActions];
                                                    newActions[index] = {
                                                        ...newActions[index],
                                                        repeat: event.target.checked,
                                                        hold: event.target.checked ? false : newActions[index].hold,
                                                    };
                                                    return newActions;
                                                });
                                            }}
                                            label="Repeat"
                                        />
                                        {action.repeat && (
                                            <NumberInput
                                                pt='sm'
                                                description="Repeat Interval"
                                                placeholder="0"
                                            />
                                        )}
                                    </Card>
                                ))}
                                <Stack align="center" justify="center" style={{flex: 1}} mr='xl'>
                                    <ActionIcon size="lg" radius="xl"
                                                onClick={() => console.log("kbd action clicked")}>
                                        <KeyboardIcon/>
                                    </ActionIcon>
                                    <ActionIcon size="lg" radius="xl"
                                                onClick={() => console.log("mouse action clicked")}>
                                        <MouseIcon/>
                                    </ActionIcon>
                                </Stack>
                            </Group>
                        </Scroller>
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

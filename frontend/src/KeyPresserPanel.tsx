import {useState} from "react";
import {Action, ActionType} from "../bindings/AutoKeyPresser";
import {ActionIcon, Card, Group, NumberInput, Radio, ScrollArea, Scroller, Stack, Switch, Text} from "@mantine/core";
import {KeyboardIcon, MouseIcon, XIcon} from "@phosphor-icons/react";
import {AddEmptyAction, DeleteAction, GetAll, UpdateAction} from "../bindings/AutoKeyPresser/keypresserservice";

const characters = [
    'Shift', 'Ctrl', 'Alt', 'Meta',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
    'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
    '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
    '`', '-', '=', '[', '\\', ';', '\'', ',', '.', '/', '~',
    '_', '+', '{', '|', ':', '"', '<', '>', '?'
];

export const KeyPresserPanel = () => {
    const [actions, setActions] = useState<Action[]>([]);

    return (
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
                        <Group>
                            <Text fw={500}
                                  mb="xs">{action.type === ActionType.ActionTypeKeyboard ? 'Keyboard' : 'Mouse'} Action {action.id}</Text>
                            <ActionIcon size="sm" radius="xl"
                                        onClick={async () => {
                                            await DeleteAction(action.id);
                                            const actions = await GetAll();
                                            if (!actions) return;
                                            setActions(actions);
                                        }}>
                                <XIcon/>
                            </ActionIcon>
                        </Group>

                        {action.type === ActionType.ActionTypeKeyboard ? (
                            <ScrollArea
                                type="always"
                                scrollbarSize={6}
                                offsetScrollbars="y"
                                h="150px"
                            >
                                <Stack>
                                    {characters.map((char, charIndex) => (
                                        <Text key={charIndex} size="sm"
                                              onClick={async () => {
                                                    await UpdateAction({
                                                        ...action,
                                                        key: char,
                                                    });
                                                    const actions = await GetAll();
                                                    if (!actions) return;
                                                    setActions(actions);
                                              }}
                                              bg={action.key === char ? 'blue.8' : undefined}
                                              py={2}
                                              px={4}
                                              style={{cursor: 'pointer', borderRadius: '4px'}}
                                        >{char}</Text>
                                    ))}
                                </Stack>
                            </ScrollArea>
                        ) : (<>
                            <Radio.Group
                                label="Mouse Button"
                                value={action.button}
                                onChange={async (value) => {
                                    await UpdateAction({
                                        ...action,
                                        button: value,
                                    });
                                    const actions = await GetAll();
                                    if (!actions) return;
                                    setActions(actions);
                                }}
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
                            onChange={async (event) => {
                                await UpdateAction({
                                    ...action,
                                    hold: event.target.checked,
                                    repeat: event.target.checked ? false : action.repeat,
                                });
                                const actions = await GetAll();
                                if (!actions) return;
                                setActions(actions);
                            }}
                            label="Hold"
                        />
                        <Switch
                            pt='md'
                            checked={action.repeat}
                            onChange={async (event) => {
                                await UpdateAction({
                                    ...action,
                                    repeat: event.target.checked,
                                    hold: event.target.checked ? false : action.hold,
                                });
                                const actions = await GetAll();
                                if (!actions) return;
                                setActions(actions);
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
                                onClick={async () => {
                                    await AddEmptyAction(ActionType.ActionTypeKeyboard);
                                    const actions = await GetAll();
                                    if (!actions) return;
                                    setActions(actions);
                                }}>
                        <KeyboardIcon/>
                    </ActionIcon>
                    <ActionIcon size="lg" radius="xl"
                                onClick={async () => {
                                    await AddEmptyAction(ActionType.ActionTypeMouse);
                                    const actions = await GetAll();
                                    if (!actions) return;
                                    setActions(actions);
                                }}>
                        <MouseIcon/>
                    </ActionIcon>
                </Stack>
            </Group>
        </Scroller>
    )
}
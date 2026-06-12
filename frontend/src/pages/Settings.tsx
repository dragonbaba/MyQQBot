import { useEffect, useState } from 'react';
import { Card } from '../components/common/Card';
import { useAppStore } from '../store/useAppStore';
import * as BotService from '../../bindings/github.com/dragonbaba/MyQQBot/internal/service/botservice.js';
import * as ConfigService from '../../bindings/github.com/dragonbaba/MyQQBot/internal/service/configservice.js';
import type { Config } from '../types';
import { Play, Square, TestTube } from 'lucide-react';

const defaultConfig: Config = {
  llm: {
    base_url: 'https://api.openai.com/v1',
    api_key: '',
    model: 'gpt-4o',
    system_prompt: '',
  },
  bot: {
    onebot_ws_url: 'ws://127.0.0.1:3001',
    admin_qq: '',
  },
  search: {
    provider: 'duckduckgo',
    tavily_api_key: '',
  },
};

const modelOptions = ['gpt-4o', 'gpt-4o-mini', 'gpt-3.5-turbo'];

export function Settings() {
  const setConfig = useAppStore((s) => s.setConfig);
  const addToast = useAppStore((s) => s.addToast);
  const [cfg, setLocalConfig] = useState<Config>(defaultConfig);
  const [activeTab, setActiveTab] = useState<'llm' | 'bot' | 'search'>('llm');
  const [starting, setStarting] = useState(false);
  const [testing, setTesting] = useState(false);

  useEffect(() => {
    ConfigService.GetConfig().then((c) => {
      setLocalConfig({
        llm: { ...defaultConfig.llm, ...c.llm },
        bot: { ...defaultConfig.bot, ...c.bot },
        search: { ...defaultConfig.search, ...c.search },
      });
      setConfig({
        llm: { ...defaultConfig.llm, ...c.llm },
        bot: { ...defaultConfig.bot, ...c.bot },
        search: { ...defaultConfig.search, ...c.search },
      });
    });
  }, [setConfig]);

  const handleSave = async () => {
    try {
      await ConfigService.SaveConfig(cfg);
      setConfig(cfg);
      addToast({ type: 'success', message: '配置已保存' });
    } catch (err) {
      addToast({ type: 'error', message: `保存失败: ${err}` });
    }
  };

  const handleTestLLM = async () => {
    setTesting(true);
    try {
      await ConfigService.TestLLMConnection(cfg.llm.base_url, cfg.llm.api_key);
      addToast({ type: 'success', message: 'LLM 连接正常' });
    } catch (err) {
      addToast({ type: 'error', message: `连接失败: ${err}` });
    } finally {
      setTesting(false);
    }
  };

  const handleToggleBot = async () => {
    setStarting(true);
    try {
      const status = await BotService.GetBotStatus();
      if (status.running) {
        await BotService.StopBot();
        addToast({ type: 'success', message: '机器人已停止' });
      } else {
        await BotService.StartBot();
        addToast({ type: 'success', message: '机器人已启动' });
      }
      const newStatus = await BotService.GetBotStatus();
      useAppStore.getState().setBotStatus(newStatus);
    } catch (err) {
      addToast({ type: 'error', message: `操作失败: ${err}` });
    } finally {
      setStarting(false);
    }
  };

  const updateLLM = (patch: Partial<Config['llm']>) =>
    setLocalConfig((prev) => ({ ...prev, llm: { ...prev.llm, ...patch } }));
  const updateBot = (patch: Partial<Config['bot']>) =>
    setLocalConfig((prev) => ({ ...prev, bot: { ...prev.bot, ...patch } }));
  const updateSearch = (patch: Partial<Config['search']>) =>
    setLocalConfig((prev) => ({ ...prev, search: { ...prev.search, ...patch } }));

  return (
    <div className="h-full overflow-auto p-6">
      <Card className="max-w-3xl mx-auto">
        <div className="flex gap-2 mb-6 border-b border-border-DEFAULT pb-2">
          {(['llm', 'bot', 'search'] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors cursor-pointer ${
                activeTab === tab
                  ? 'bg-brand-cta/15 text-brand-cta'
                  : 'text-text-secondary hover:text-text-primary hover:bg-surface-elevated'
              }`}
            >
              {tab === 'llm' && 'LLM 配置'}
              {tab === 'bot' && 'QQ 机器人'}
              {tab === 'search' && '搜索设置'}
            </button>
          ))}
        </div>

        <div className="space-y-4">
          {activeTab === 'llm' && (
            <>
              <label className="block text-sm text-text-secondary">
                Base URL
                <input
                  value={cfg.llm.base_url}
                  onChange={(e) => updateLLM({ base_url: e.target.value })}
                  placeholder="https://api.openai.com/v1"
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>
              <label className="block text-sm text-text-secondary">
                API Key
                <input
                  type="password"
                  value={cfg.llm.api_key}
                  onChange={(e) => updateLLM({ api_key: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>
              <label className="block text-sm text-text-secondary">
                Model
                <select
                  value={cfg.llm.model}
                  onChange={(e) => updateLLM({ model: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                >
                  {modelOptions.map((m) => (
                    <option key={m} value={m}>
                      {m}
                    </option>
                  ))}
                </select>
              </label>
              <label className="block text-sm text-text-secondary">
                System Prompt
                <textarea
                  value={cfg.llm.system_prompt}
                  onChange={(e) => updateLLM({ system_prompt: e.target.value })}
                  rows={4}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50 resize-none"
                />
              </label>
              <button
                onClick={handleTestLLM}
                disabled={testing}
                className="inline-flex items-center gap-2 px-4 py-2 bg-status-info/15 text-status-info hover:bg-status-info/25 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
              >
                <TestTube className="w-4 h-4" />
                {testing ? '测试中...' : '测试连接'}
              </button>
            </>
          )}

          {activeTab === 'bot' && (
            <>
              <label className="block text-sm text-text-secondary">
                OneBot WebSocket URL
                <input
                  value={cfg.bot.onebot_ws_url}
                  onChange={(e) => updateBot({ onebot_ws_url: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>
              <label className="block text-sm text-text-secondary">
                Admin QQ
                <input
                  value={cfg.bot.admin_qq}
                  onChange={(e) => updateBot({ admin_qq: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>
              <button
                onClick={handleToggleBot}
                disabled={starting}
                className="inline-flex items-center gap-2 px-4 py-2 bg-brand-cta/15 text-brand-cta hover:bg-brand-cta/25 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
              >
                {useAppStore((s) => s.botStatus.running) ? (
                  <>
                    <Square className="w-4 h-4" /> 停止机器人
                  </>
                ) : (
                  <>
                    <Play className="w-4 h-4" /> 启动机器人
                  </>
                )}
              </button>
            </>
          )}

          {activeTab === 'search' && (
            <>
              <div className="space-y-2">
                <p className="text-sm text-text-secondary">搜索 Provider</p>
                <label className="flex items-center gap-2 text-sm text-text-primary cursor-pointer">
                  <input
                    type="radio"
                    name="provider"
                    checked={cfg.search.provider === 'duckduckgo'}
                    onChange={() => updateSearch({ provider: 'duckduckgo' })}
                    className="cursor-pointer"
                  />
                  DuckDuckGo（免费）
                </label>
                <label className="flex items-center gap-2 text-sm text-text-primary cursor-pointer">
                  <input
                    type="radio"
                    name="provider"
                    checked={cfg.search.provider === 'tavily'}
                    onChange={() => updateSearch({ provider: 'tavily' })}
                    className="cursor-pointer"
                  />
                  Tavily
                </label>
              </div>
              {cfg.search.provider === 'tavily' && (
                <label className="block text-sm text-text-secondary">
                  Tavily API Key
                  <input
                    type="password"
                    value={cfg.search.tavily_api_key}
                    onChange={(e) => updateSearch({ tavily_api_key: e.target.value })}
                    className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                  />
                </label>
              )}
            </>
          )}
        </div>

        <div className="mt-6 flex justify-end">
          <button
            onClick={handleSave}
            className="px-6 py-2 bg-brand-cta hover:bg-brand-cta/80 text-surface-base font-medium rounded-lg transition-colors cursor-pointer"
          >
            保存配置
          </button>
        </div>
      </Card>
    </div>
  );
}

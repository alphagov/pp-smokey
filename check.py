#!/usr/bin/env python
# encoding: utf-8

import sys
import json
from collections import namedtuple
from operator import add


Step = namedtuple('Step', 'name, status')
Scenario = namedtuple('Scenario', 'name, steps')
Feature = namedtuple('Feature', 'name, uri, scenarios')


def main():
    feature_name, json_file = parse_args(sys.argv)

    feature = get_feature(load_json(json_file), feature_name)
    log_result_and_exit(feature)


# Utils
def ascii(value):
    return value.encode('ascii', 'ignore')


def parse_args(args):
    return args[1], args[2]


def load_json(json_file):
    with open(json_file) as f:
        return json.load(f)


# Loading from JSON
def get_feature(smokey_json, feature_name):
    for feature_json in smokey_json:
        if feature_json['uri'] == get_feature_uri(feature_name):
            return Feature(
                ascii(feature_json['name']),
                ascii(feature_json['uri']),
                map(get_scenario, find_scenarios(feature_json)))


def get_feature_uri(feature_name):
    return 'features/{}.feature'.format(feature_name)


def find_scenarios(feature_json):
    return [element for element in feature_json['elements']
            if element['keyword'] == 'Scenario']


def get_scenario(scenario_json):
    return Scenario(
        ascii(scenario_json['name']),
        map(get_step, scenario_json['steps']))


def get_step(step):
    return Step(
        "{}{}".format(ascii(step['keyword']), ascii(step['name'])),
        step['result']['status'])


# Counting steps
def count_failing_steps(feature):
    return count_steps_by_status(feature, 'failed')


def count_passing_steps(feature):
    return count_steps_by_status(feature, 'passed')


def count_steps_by_status(feature, status):
    return len([step for step in get_feature_steps(feature)
               if step.status == status])


def get_feature_steps(feature):
    return reduce(add,
                  [scenario.steps for scenario in feature.scenarios])


# Rendering as text
def feature_message(feature):
    return ('{failing} failed, {passing} passed;\n'
            '{name} ({uri})\n{scenarios}').format(
                failing=count_failing_steps(feature),
                passing=count_passing_steps(feature),
                name=feature.name,
                uri=feature.uri,
                scenarios="\n".join(map(scenario_message, feature[2])))


def scenario_message(scenario):
    return "  Scenario: {name}\n{steps}".format(
        name=scenario.name,
        steps="\n".join(map(step_message, scenario.steps)))


def step_message(step):
    return '    Step: [{status}] {name}'.format(
        name=step.name,
        status='PASS' if step.status == "passed" else 'FAIL')


# Status message for Sensu
def log_result_and_exit(feature):
    exit_status = 2 if count_failing_steps(feature) > 0 else 0
    message = feature_message(feature)

    print(message)
    sys.exit(exit_status)

if __name__ == '__main__':
    main()

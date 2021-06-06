cmake_minimum_required(VERSION 2.8.12)

include_directories (${CMAKE_CURRENT_SOURCE_DIR}/include)
set(RNX2RTKP_SRC_DIR ${CMAKE_CURRENT_SOURCE_DIR}/src)

add_library(rtklib STATIC
				${RNX2RTKP_SRC_DIR}/datum.c
				${RNX2RTKP_SRC_DIR}/download.c
				${RNX2RTKP_SRC_DIR}/rtkcmn.c
				${RNX2RTKP_SRC_DIR}/rinex.c
				${RNX2RTKP_SRC_DIR}/rtkpos.c
				${RNX2RTKP_SRC_DIR}/postpos.c
				${RNX2RTKP_SRC_DIR}/solution.c
				${RNX2RTKP_SRC_DIR}/lambda.c
				${RNX2RTKP_SRC_DIR}/geoid.c
				${RNX2RTKP_SRC_DIR}/sbas.c
				${RNX2RTKP_SRC_DIR}/preceph.c
				${RNX2RTKP_SRC_DIR}/pntpos.c
				${RNX2RTKP_SRC_DIR}/ephemeris.c
				${RNX2RTKP_SRC_DIR}/options.c
				${RNX2RTKP_SRC_DIR}/ppp.c
				${RNX2RTKP_SRC_DIR}/ppp_ar.c
				${RNX2RTKP_SRC_DIR}/rtcm.c
				${RNX2RTKP_SRC_DIR}/rtcm2.c
				${RNX2RTKP_SRC_DIR}/rtcm3.c
				${RNX2RTKP_SRC_DIR}/rtcm3e.c
				${RNX2RTKP_SRC_DIR}/ionex.c
				${RNX2RTKP_SRC_DIR}/qzslex.c
				${RNX2RTKP_SRC_DIR}/convrnx.c
				${RNX2RTKP_SRC_DIR}/tle.c)

SOURCE_GROUP("src" FILES ${POINTS2GRID_CPP})
SOURCE_GROUP("include" FILES ${POINTS2GRID_HPP})

IF(FLYWAVE_ENABLE_SOLUTION_FOLDERS)
	SET_TARGET_PROPERTIES(rtklib PROPERTIES FOLDER external)
ENDIF(FLYWAVE_ENABLE_SOLUTION_FOLDERS)
			   